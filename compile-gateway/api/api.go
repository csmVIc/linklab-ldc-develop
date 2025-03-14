package api

import (
	"errors"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/compile-gateway/api/handler"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ServerAddress HTTP服务绑定的地址
type ServerAddress struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// AInfo HTTP服务的配置参数
type AInfo struct {
	Address ServerAddress `json:"address"`
	Handler handler.HInfo `json:"handler"`
}

// RouterInit 路由初始化
func RouterInit(r *gin.Engine, i *handler.HInfo) error {

	if i == nil {
		err := errors.New("api.RouterInit input IInfo is nil")
		return err
	}

	// 用户验证实例化
	var err error
	auth.Adriver, err = auth.New("users")
	if err != nil {
		log.Errorf("auth.New error {%v}", err)
		return err
	}

	group := r.Group("/api")
	group.Use(auth.Adriver.AuthorizationRequired())

	if err := handler.RouterInit(group, i); err != nil {
		log.Errorf("api.RouterInit compile init return error {%v}", err)
		return err
	}
	return nil
}
