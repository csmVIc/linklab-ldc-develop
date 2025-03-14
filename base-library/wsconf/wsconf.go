package wsconf

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	// UpgraderGlobal 全局实例
	UpgraderGlobal = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
