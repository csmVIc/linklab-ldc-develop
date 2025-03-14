package eapi

import (
	"errors"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/client/iotnode/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	einfo *EInfo
)

// RouterInit 路由初始化
func RouterInit(r *gin.Engine, ei *EInfo) error {

	if ei == nil {
		err := errors.New("RouterInit ei is nil")
		log.Error(err)
		return err
	}
	einfo = ei

	// 服务端验证
	aRouterGroup := r.Group("/api")
	aRouterGroup.Use(auth.EdgeClientAuthorizationRequired(api.ADriver.GetToken()))

	aRouterGroup.GET("/echo", echohandler)
	aRouterGroup.POST("/pod/apply", podapplyhandler)
	aRouterGroup.GET("/pod/log", podloghandler)
	aRouterGroup.POST("/pod/delete", poddeletehandler)
	aRouterGroup.GET("/pod/exec", podexechandler)
	aRouterGroup.POST("/image/build", imagebuildhandler)

	return nil
}

func echohandler(c *gin.Context) {
	echo(c)
}

func podapplyhandler(c *gin.Context) {
	podapply(c)
}

func podloghandler(c *gin.Context) {
	podlog(c)
}

func poddeletehandler(c *gin.Context) {
	poddelete(c)
}

func podexechandler(c *gin.Context) {
	podexec(c)
}

func imagebuildhandler(c *gin.Context) {
	imagebuild(c)
}
