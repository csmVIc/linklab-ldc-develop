package api

import (
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/user-service/api/client"
	"linklab/device-control-v2/user-service/api/device"
	"linklab/device-control-v2/user-service/api/edgenode"
	"linklab/device-control-v2/user-service/api/ws"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RouterInit 路由初始化
func RouterInit(r *gin.Engine, ai *AInfo) error {
	// 允许所有访问
	// r.Use(cors.Default())

	// 用户验证实例化
	var err error
	auth.Adriver, err = auth.New("users")
	if err != nil {
		log.Errorf("auth.New error {%v}", err)
		return err
	}

	// /api路径下的接口主要和设备操作有关，均需要用户验证
	aRouterGroup := r.Group("/api")
	aRouterGroup.Use(auth.Adriver.AuthorizationRequired())

	if err := device.RouterInit(aRouterGroup, &ai.Device); err != nil {
		log.Errorf("device.RouterInit error {%v}", err)
		return err
	}

	if err := ws.RouterInit(aRouterGroup, &ai.Ws); err != nil {
		log.Errorf("ws.RouterInit error {%v}", err)
		return err
	}

	if err := client.RouterInit(aRouterGroup, &ai.Client); err != nil {
		log.Errorf("client.RouterInit error {%v}", err)
		return err
	}

	if err := edgenode.RouterInit(aRouterGroup, &ai.EdgeNode); err != nil {
		log.Errorf("edgenode.RouterInit error {%v}", err)
		return err
	}

	return nil
}
