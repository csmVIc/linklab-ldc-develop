package ws

import (
	"errors"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/wsconf"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func wshandler(c *gin.Context) {

	// 获取用户id
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// websocket 建立回复
	protocols := websocket.Subprotocols(c.Request)
	var httpHeader http.Header = nil
	if len(protocols) > 0 {
		log.Debugf("websocket subprotocols {%v}", protocols)
		httpHeader = http.Header{
			"Sec-Websocket-Protocol": []string{protocols[0]},
		}
	}

	// 建立websocket
	if handler, err := wsconf.UpgraderGlobal.Upgrade(c.Writer, c.Request, httpHeader); err != nil {
		log.Errorf("websocket upgrade error {%v}", err)
	} else {
		log.Info("websocket upgrade success")
		monitor(handler, userid, winfo.TimeOut, winfo.ChanSize, winfo.Msg.Topic)
	}
}
