package edgenode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func podLogMsgForward(userhandler *websocket.Conn, edgehandler *websocket.Conn) {

	edgeExistChan := make(chan error, 1)
	msgForwardChan := make(chan *response.EdgeClientPodLog, einfo.PodLog.MsgForwardChanSize)

	go func() {
		for {
			edgeResp := &response.EdgeClientPodLog{}
			if err := edgehandler.ReadJSON(&edgeResp); err != nil {
				err = fmt.Errorf("edgehandler.ReadJSON error {%v}", err)
				log.Error(err)
				edgeExistChan <- err
				return
			}
			msgForwardChan <- edgeResp
			log.Debugf("read json {%v} {%v}", edgeResp.Type, edgeResp.Msg)
		}
	}()

	for EXIST := false; !EXIST; {
		select {
		case edgeResp := <-msgForwardChan:
			if edgeResp.Type == response.NormalPodLog {
				if err := userhandler.WriteJSON(response.Response{Code: 0, Msg: edgeResp.Msg}); err != nil {
					log.Errorf("userhandler.WriteJSON error {%v}", err)
					EXIST = true
					break
				}
			} else if edgeResp.Type == response.ErrorPodLog {
				if err := userhandler.WriteJSON(response.Response{Code: -1, Msg: edgeResp.Msg}); err != nil {
					log.Errorf("userhandler.WriteJSON error {%v}", err)
				}
				EXIST = true
				break
			}

			log.Debugf("send json {%v} {%v}", edgeResp.Type, edgeResp.Msg)

		case err := <-edgeExistChan:
			log.Errorf("edgeExistChan error {%v}", err)
			EXIST = true
			break

		case <-time.After(time.Second):
			if err := userhandler.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("userhandler.WriteMessage ping error {%v}", err)
				EXIST = true
				break
			}
		}
	}

	userhandler.Close()
	edgehandler.Close()
}
