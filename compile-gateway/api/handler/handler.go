package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	hinfo *HInfo
)

// RouterInit compile路由初始化
func RouterInit(r *gin.RouterGroup, hi *HInfo) error {

	if hi == nil {
		err := errors.New("RouterInit hi is nil")
		log.Error(err)
		return err
	}

	hinfo = hi

	r.POST("/compile", posthandler)
	r.POST("/compilesystem", compilesystemhandler)
	r.GET("/compile/nonblock", getnonblockhandler)
	r.GET("/compile/block", getblockhandler)
	r.GET("/compile/block/status", getblockstatushandler)
	return nil
}

func posthandler(c *gin.Context) {
	Uploadhandler(c)
}

func compilesystemhandler(c *gin.Context) {
	CompileSystem(c)
}

func getnonblockhandler(c *gin.Context) {
	DownloadNonBlockhandler(c)
}

func getblockhandler(c *gin.Context) {
	DownloadBlockhandler(c)
}

func getblockstatushandler(c *gin.Context) {
	GetBlockStatusHandler(c)
}
