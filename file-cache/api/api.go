package api

import (
	"linklab/device-control-v2/file-cache/api/fcache"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RouterInit 路由初始化
func RouterInit(r *gin.Engine, fci *fcache.FCInfo) error {

	if err := fcache.RouterInit(r, fci); err != nil {
		log.Errorf("fcache.RouterInit error {%v}", err)
		return err
	}

	return nil
}
