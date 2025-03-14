package edgenode

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/client/iotnode/api"
    "net/http"  // 添加
    "time"      // 添加
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
)

// ImageBuild 镜像构建
func (ed *Driver) ImageBuild(imageName string, fileHash string, nameSpace string, nodeselector map[string]string) error {

	// 读取yaml文件
	log.Debugf("准备读取构建模板文件: ./yaml/build.yaml")
	buildyaml, err := ioutil.ReadFile("./yaml/build.yaml")
	if err != nil {
		err = fmt.Errorf("ioutil.ReadFile error {%v}", err)
		log.Error(err)
		return err
	}
	log.Debugf("构建模板文件读取成功, 大小: %d bytes", len(buildyaml))

	// 环境变量
	env := []apiv1.EnvVar{}
	// 设置镜像名称
	imageName = fmt.Sprintf("%v/%v/%v", ed.info.ImageBuild.RegistryAddress, nameSpace, imageName)
	log.Debugf("设置完整镜像名称: %s", imageName)
	env = append(env, apiv1.EnvVar{
		Name:  "IMAGE_NAME",
		Value: imageName,
	})

	// 设置token
	token := api.ADriver.GetToken()
	// log.Debugf("获取到Token: %s", token)
	env = append(env, apiv1.EnvVar{
		Name:  "TOKEN",
		Value: token,
	})

	// 设置文件Hash
	log.Debugf("设置文件Hash环境变量: %s", fileHash)
	env = append(env, apiv1.EnvVar{
		Name:  "FILE_HASH",
		Value: fileHash,
	})
   // 设置下载URL
   log.Debugf("设置构建下载URL: %s", ed.info.ImageBuild.BuildDownloadURL)
	env = append(env, apiv1.EnvVar{
		Name:  "BUILD_DOWNLOAD_URL",
		Value: ed.info.ImageBuild.BuildDownloadURL,
	})
	log.Debugf("环境变量设置完成, 总计 %d 个变量", len(env))


	// 在创建 Pod 之前添加网络连通性检查
	log.Debugf("正在检查目标服务可用性...")
	timeout := time.Second * 5
	client := http.Client{
		Timeout: timeout,
	}
	testURL := fmt.Sprintf("%s/api/imagebuild?filehash=%s", ed.info.ImageBuild.BuildDownloadURL, fileHash)
	resp, err := client.Get(testURL)
	if err != nil {
		log.Warnf("预检查发现目标服务可能不可用: %v", err)
	} else {
		resp.Body.Close()
		log.Debugf("目标服务连通性检查通过")
	}
	
	// 创建pod
	if err := ed.PodApply(buildyaml, nameSpace, env, false, nodeselector, ""); err != nil {
		// err = fmt.Errorf("%v", err)
		// log.Error(err)
		log.Errorf("创建构建Pod失败: %v", err)
		return err
	}
	log.Debugf("构建Pod创建成功")


	return nil
}
