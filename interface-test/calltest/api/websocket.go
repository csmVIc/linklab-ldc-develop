package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// GetWebSocketHandle 获得websocket句柄
func (ad *Driver) GetWebSocketHandle(token string) (*websocket.Conn, error) {

	handle, resp, err := websocket.DefaultDialer.Dial(ad.info.WebSocket.URL,
		http.Header{"Authorization": []string{token}})
	if err != nil {
		if resp != nil {
			respbody := make([]byte, resp.ContentLength)
			if _, readerr := resp.Body.Read(respbody); readerr != nil {
				log.Errorf("resp.Body.Read error {%v}", readerr)
				return nil, readerr
			}
			err = fmt.Errorf("websocket.DefaultDialer.Dial error {%v} resp {%v}", err, string(respbody))
		} else {
			err = fmt.Errorf("websocket.DefaultDialer.Dial error {%v}", err)
		}
		log.Error(err)
		return nil, err
	}

	return handle, nil
}
