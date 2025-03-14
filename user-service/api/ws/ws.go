package ws

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	winfo *WInfo
)

// RouterInit websocket路由初始化
func RouterInit(r *gin.RouterGroup, wi *WInfo) error {
	if wi == nil {
		err := errors.New("RouterInit wi is nil")
		log.Error(err)
		return err
	}
	winfo = wi
	// websocket处理句柄
	r.GET("/ws", wshandler)

	// websocket日志过滤句柄
	r.GET("/wsfilter", wsfilterhandler)
	return nil
}
