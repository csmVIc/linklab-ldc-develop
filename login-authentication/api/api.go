package api

import (
	"linklab/device-control-v2/login-authentication/api/user"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// RouterInit 路由初始化
func RouterInit(r *gin.Engine, ui *user.UInfo) error {

	// /user路径下的接口和用户有关，部分需要用户验证
	uRouterGroup := r.Group("/user")
	if err := user.RouterInit(uRouterGroup, ui); err != nil {
		log.Errorf("user.RouterInit error {%v}", err)
		return err
	}
	// /client路径下的接口和客户端有关,部分需要客户端验证
	// cRouterGroup := r.Group("/client")
	// if err := client.RouterInit(cRouterGroup, ci); err != nil {
	// 	log.Errorf("client.RouterInit error {%v}", err)
	// 	return err
	// }
	return nil
}
