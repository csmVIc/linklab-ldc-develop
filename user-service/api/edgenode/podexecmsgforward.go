package edgenode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/parameter/response"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func podExecMsgForward(userhandler *websocket.Conn, edgehandler *websocket.Conn) {

	defer func() {
		userhandler.Close()
		edgehandler.Close()
	}()

	userExistChan := make(chan error, 1)
	edgeExistChan := make(chan error, 1)

	// 用户消息转发到边缘
	go func() {
		for {
			execMsg := &msg.EdgeClientPodExec{}
			if err := userhandler.ReadJSON(execMsg); err != nil {
				err = fmt.Errorf("userhandler.ReadJSON error {%v}", err)
				log.Error(err)
				userExistChan <- err
				return
			}
			if err := edgehandler.WriteJSON(execMsg); err != nil {
				err = fmt.Errorf("edgehandler.WriteJSON error {%v}", err)
				log.Error(err)
				userExistChan <- err
				return
			}
			log.Debugf("edgehandler.WriteJSON {%v}", execMsg.Msg)
			if execMsg.Type == msg.ErrorPodExec {
				err := fmt.Errorf("userhandler.ReadJSON error {%v}", execMsg.Msg)
				log.Error(err)
				userExistChan <- err
				return
			}
		}
	}()

	// 边缘消息转发到用户
	go func() {
		for {
			execMsg := &msg.EdgeClientPodExec{}
			if err := edgehandler.ReadJSON(execMsg); err != nil {
				err = fmt.Errorf("edgehandler.ReadJSON error {%v}", err)
				log.Error(err)
				edgeExistChan <- err
				return
			}

			resp := response.Response{
				Code: 0,
				Msg:  execMsg.Msg,
			}
			if execMsg.Type == msg.ErrorPodExec {
				resp.Code = -1
			}
			if err := userhandler.WriteJSON(resp); err != nil {
				err = fmt.Errorf("userhandler.WriteJSON error {%v}", err)
				log.Error(err)
				edgeExistChan <- err
				return
			}
			log.Debugf("userhandler.WriteJSON {%v}", execMsg.Msg)
			if execMsg.Type == msg.ErrorPodExec {
				err := fmt.Errorf("edgehandler.ReadJSON error {%v}", execMsg.Msg)
				log.Error(err)
				edgeExistChan <- err
				return
			}
		}
	}()

	// 退出检测
	for EXIST := false; !EXIST; {
		select {
		case err := <-userExistChan:
			log.Errorf("userExistChan error {%v}", err)
			EXIST = true
			break

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
			if err := edgehandler.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("edgehandler.WriteMessage ping error {%v}", err)
				EXIST = true
				break
			}
		}
	}
}
