package eapi

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func imagebuild(c *gin.Context) {
	log.Debugf("收到镜像构建请求:")
	log.Debugf("- 请求头Authorization: %s", c.GetHeader("Authorization"))
	log.Debugf("- 请求类型Content-Type: %s", c.ContentType())

	// 参数验证
	var p request.EdgeClientImageBuild
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	log.Debugf("解析到的请求参数:")
	log.Debugf("- 镜像名称: %s", p.ImageName)
	log.Debugf("- 文件Hash: %s", p.FileHash)
	log.Debugf("- 命名空间: %s", p.Namespace)
	log.Debugf("- 节点选择器: %+v", p.NodeSelector)
	// 创建命名空间
	if err := edgenode.EDriver.NamespaceCreateIfNotExist(p.Namespace); err != nil {
		err = fmt.Errorf("namespace create error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	log.Debugf("命名空间创建/确认成功: %s", p.Namespace)

	// 镜像构建
	log.Debugf("开始构建镜像:")
	log.Debugf("- 镜像名称: %s", p.ImageName)
	log.Debugf("- 使用的文件Hash: %s", p.FileHash)
	log.Debugf("- 在命名空间: %s", p.Namespace)
	// 镜像构建
	if err := edgenode.EDriver.ImageBuild(p.ImageName, p.FileHash, p.Namespace, p.NodeSelector); err != nil {
		// err = fmt.Errorf("%v", err)
		// log.Error(err)
		log.Errorf("镜像构建失败: %v", err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}   
	log.Debugf("镜像构建成功: %s", p.ImageName)


	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
}
