package eapi

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
)


func podapply(c *gin.Context) {
	// 添加请求调试信息
	log.Debugf("收到Pod部署请求:")
	log.Debugf("- Authorization: %s", c.GetHeader("Authorization"))
	log.Debugf("- Content-Type: %s", c.ContentType())

	// 参数验证
	var p request.EdgeClientPodApply
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
    log.Debugf("请求参数: %+v", p)

	// 下载PodYaml文件
	podyaml, err := api.ADriver.PodYamlDownload(p.YamlHash)
	if err != nil {
		err = fmt.Errorf("pod yaml {%v} not exist error {%v}", p.YamlHash, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 创建命名空间
	if err := edgenode.EDriver.NamespaceCreateIfNotExist(p.Namespace); err != nil {
		err = fmt.Errorf("namespace create error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// 先创建cloud ------
	log.Debugf("先创建 自制文件 cloud - pod")
	// 如果创建的是pod是cloud，继续创建第二个pod - edge
	// 一. 先对pod.yaml修改
	// 1. 修改pod名称,labels：metadata.name = edge-test，metadata.labels.app = edge-test
	// 2. 修改container名称, 镜像：spec.containers[0].name = edge-test，spec.containers[0].image = edge_test:v2.0
	// 3. 增加一个环境变量，spec.containers[0].env = ENABLE_EDGE_SOCKET=true
	// 4. 删除args配置
	// 创建edge pod的YAML配置
cloudPodYaml := []byte(`apiVersion: v1
apiVersion: v1
kind: Pod
metadata:
  name: cloud-test
  labels:
    app: cloud-test
spec:
  restartPolicy: OnFailure
  containers:
  - name: cloud-test
    image: edge_cloud:v1.0
    env:
    - name: PYTHONUNBUFFERED
      value: "1"
    args: ["/bin/bash","-c","python3 cloud_socket.py -i 0.0.0.0 -p 99"]
    ports:
    - containerPort: 99
      protocol: TCP

`)

	log.Debugf("创建了新的cloud pod YAML配置")
	// 二. 创建ingress
	log.Debugf("开始配置Ingress - cloud")
	var ingressMap *map[string]string
	if ingressMap, err = edgenode.EDriver.LinkPodIngress(cloudPodYaml, p.Namespace); err != nil {
		// err = fmt.Errorf("%v", err)
		log.Errorf("cloud - Ingress配置失败: %v", err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// 三. 创建env
	env := []apiv1.EnvVar{}
	if ingressMap != nil && len(*ingressMap) > 0 {
		log.Debugf("处理cloud - Ingress环境变量")
		for key, value := range *ingressMap {
			env = append(env, apiv1.EnvVar{
				Name:  fmt.Sprintf("INGRESS_PORT%s", key),
				Value: value,
			})
		}
	}
	// 四. 创建pod
	log.Debugf("节点选择器: %+v", p.NodeAddSelector)
	if err := edgenode.EDriver.PodApply(cloudPodYaml, p.Namespace, env, p.UseEdgeRegistry, nil, p.NodeAddSelector); err != nil {
		// err = fmt.Errorf("%v", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	log.Debugf("自制文件 cloud - Pod部署成功")

	// 创建完成 ------
	// 关联Ingress
	ingressMap = nil
	p.CreateIngress = true
	if p.CreateIngress {
		log.Debugf("开始配置Ingress")
		// 添加YAML解析日志
		log.Debugf("Pod YAML配置详情:")

		if ingressMap, err = edgenode.EDriver.LinkPodIngress(podyaml, p.Namespace); err != nil {
			// err = fmt.Errorf("%v", err)
			log.Errorf("Ingress配置失败: %v", err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		if ingressMap == nil {
			log.Warnf("Ingress配置成功但未返回映射信息 ingressMap == nil")
		} else {
			// 创建成功
			log.Debugf("Ingress配置成功，映射详情:")
			for key, value := range *ingressMap {
				log.Debugf("- 端口映射 %s -> %s", key, value)
			}
		}
	}

	env = []apiv1.EnvVar{}
	if ingressMap != nil && len(*ingressMap) > 0 {
		log.Debugf("处理Ingress环境变量")
		for key, value := range *ingressMap {
			env = append(env, apiv1.EnvVar{
				Name:  fmt.Sprintf("INGRESS_PORT%s", key),
				Value: value,
			})
			log.Debugf("添加环境变量: INGRESS_PORT%s = %s", key, value)
		}
		log.Debugf("环境变量处理完成, 总数: %d", len(env))
	}

	// 部署Pod
    log.Debugf("开始部署Pod到命名空间: %s", p.Namespace)
    log.Debugf("使用Edge Registry: %v", p.UseEdgeRegistry)
    log.Debugf("节点选择器: %+v", p.NodeAddSelector)
	if err := edgenode.EDriver.PodApply(podyaml, p.Namespace, env, p.UseEdgeRegistry, nil, p.NodeAddSelector); err != nil {
		// err = fmt.Errorf("%v", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	log.Debugf("Pod部署成功")


	if ingressMap != nil {
		log.Debugf("返回带Ingress映射的成功响应")
		c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: response.EdgeClientPodApply{
			IngressMap: *ingressMap,
		}})
		return
	} else {
		log.Debugf("返回普通成功响应")
		c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
		return
	}
}