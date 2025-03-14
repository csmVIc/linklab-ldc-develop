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

func poddelete(c *gin.Context) {

	// 参数验证
	p := request.EdgeClientPodDelete{}
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 删除服务
	if err := edgenode.EDriver.ServiceDelete(p.Namespace, p.Pod); err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 删除Ingress
	if err := edgenode.EDriver.IngressDelete(p.Namespace, p.Pod); err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 删除Pod
	if err := edgenode.EDriver.PodDelete(p.Namespace, p.Pod); err != nil {
		// err = fmt.Errorf("pod delete error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
}
