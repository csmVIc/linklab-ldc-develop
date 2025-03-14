package ws

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func monitor(handler *websocket.Conn, userid string, timeout int, chansize int, topic string) {

	log.Debugf("user {%v} monitor websocket create", userid)
	// 退出时需要断开websocket连接
	defer func() {
		if err := handler.Close(); err != nil {
			log.Errorf("user {%v} close weboscket error {%v}", userid, err)
		}
		log.Debugf("user {%v} monitor websocket destroy", userid)
	}()

	// 创建消息转发channel,以及取消订阅channel
	msgchan, exitchan := make(chan *nats.Msg, chansize), make(chan bool)
	defer func() {
		exitchan <- true
		close(msgchan)
		close(exitchan)
	}()
	go func() {
		for true {
			err := subscriber(&msgchan, &exitchan, fmt.Sprintf(topic, userid))
			if err != nil {
				log.Errorf("user {%v} subscriber error {%v}", userid, err)
				time.Sleep(time.Second)
				log.Infof("user {%v} subscriber reconnect", userid)
				continue
			} else {
				log.Infof("user {%v} subscriber exit", userid)
				break
			}
		}
	}()

	// 转发消息,定期检测websocket连接状况,以及消息队列的连接状况
	for count := 0; count < timeout; {
		select {
		case natmsg := <-msgchan:

			// 增加时间戳
			uMsg := msg.UserMsg{}
			if err := json.Unmarshal(natmsg.Data, &uMsg); err != nil {
				log.Errorf("user {%v} json.Unmarshal error {%v}", userid, err)
				return
			}
			uMsg.TimeStamp = time.Now().UnixNano()

			// JSON序列化
			binary, err := json.Marshal(uMsg)
			if err != nil {
				log.Errorf("user {%v} json.Marshal error {%v}", userid, err)
				return
			}

			if err := handler.WriteMessage(websocket.TextMessage, binary); err != nil {
				log.Errorf("user {%v} websocket write msg error {%v}", userid, err)
				return
			}

			// 调试
			log.Debugf("userid {%v} websocket write msg {%v}", userid, string(binary))

			count = 0
		case <-time.After(time.Second):
			count++
			if err := handler.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Errorf("user {%v} websocket write ping msg error {%v}", userid, err)
				return
			}
		}
	}

	// websocket很长一段时间没有接收到消息,选择主动断开连接
	err := fmt.Errorf("user {%v} websocket has not sent a message for a long time {%vs}, so it is disconnected", userid, timeout)
	log.Error(err)
	cMsg := msg.UserMsg{
		Code:      -1,
		Type:      msg.UserControlMsg,
		TimeStamp: time.Now().UnixNano(),
		Data: msg.UserControlData{
			Msg: err.Error(),
		},
	}
	if err := handler.WriteJSON(cMsg); err != nil {
		log.Errorf("user {%v} websocket writejson error {%v}", userid, err)
	}
}
