package client

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	cinfo *CInfo
)

// RouterInit 设备操作路由初始化
func RouterInit(r *gin.RouterGroup, ci *CInfo) error {

	if ci == nil {
		err := errors.New("RouterInit di is nil")
		log.Error(err)
		return err
	}
	cinfo = ci

	r.GET("/client", listclienthandler)
	r.POST("/client/tenant/change", changeclienttenanthandler)
	r.POST("/client/tenant/add", addclienttenanthandler)
	r.POST("/client/tenant/sub", subclienttenanthandler)

	return nil
}

func listclienthandler(c *gin.Context) {
	listclient(c)
}

func changeclienttenanthandler(c *gin.Context) {
	changeclienttenant(c)
}

func addclienttenanthandler(c *gin.Context) {
	addclienttenant(c)
}

func subclienttenanthandler(c *gin.Context) {
	subclienttenant(c)
}
