package msgforward

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func (md *Driver) cmdhandler(subchansize int, msgtopic string, mqtttopic string) {

	defer func() {
		md.exitsignal = true
	}()

	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		log.Errorf("get nats conn error {%v}", err)
		return
	}

	submsgchan := make(chan *nats.Msg, subchansize)
	defer func() {
		close(submsgchan)
	}()
	sub, err := natsconn.QueueSubscribe(msgtopic, "device-service", func(msg *nats.Msg) {
		submsgchan <- msg
		log.Infof("sub topic {%v} msg {%v}", msg.Subject, string(msg.Data))
	})
	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Errorf("unsubscribe messenger topic {%v} error {%v}", sub.Subject, err)
		}
	}()

	for {
		select {
		case nmsg := <-submsgchan:
			username, err := messenger.GetIDFromTopic(nmsg.Subject, "clients")
			if err != nil {
				log.Errorf("messenger.GetIDFromTopic error {%v}", err)
				return
			}
			clientid, err := auth.GetClientIDByUserName(username)
			if err != nil {
				log.Errorf("auth.GetClientIDByUserName error {%v}", err)
				return
			}

			// 转发至设备管理客户端
			err = mqttclient.MDriver.PubMsgByte(fmt.Sprintf(mqtttopic, username, clientid), 2, nmsg.Data)
			if err != nil {
				log.Errorf("mqttclient.MDriver.PubMsgByte error {%v}", err)
				return
			}
			log.Debugf("mqtt pub msg {%v} {%v}", fmt.Sprintf(mqtttopic, username, clientid), string(nmsg.Data))

			// 回复消息
			respond := msg.ReplyMsg{
				Code: 0,
				Msg:  "success",
				Data: nil,
			}
			respondBytes, err := json.Marshal(respond)
			if err != nil {
				log.Errorf("json.Marshal error {%v}", err)
				return
			}
			if err := nmsg.Respond(respondBytes); err != nil {
				log.Errorf("nmsg.Respond error {%v}", err)
				return
			}

			// 记录命令日志
			devicecmd := msg.DeviceCmd{}
			if err := json.Unmarshal(nmsg.Data, &devicecmd); err != nil {
				log.Errorf("json.Unmarshal error {%v}", err)
				return
			}
			tags := map[string]string{
				"groupid": devicecmd.GroupID,
			}
			fields := map[string]interface{}{
				"taskindex": devicecmd.TaskIndex,
				"clientid":  username,
				"deviceid":  devicecmd.DeviceID,
				"cmd":       devicecmd.Cmd,
				"podname":   os.Getenv("POD_NAME"),
				"nodename":  os.Getenv("NODE_NAME"),
			}
			err = logger.Ldriver.WriteLog("cmdwrite", tags, fields)
			if err != nil {
				log.Errorf("write log err {%v}", err)
				return
			}

		case <-time.After(time.Second):
			if messenger.Mdriver.GetClosed() == true {
				log.Errorf("messenger.Mdriver.GetClosed true error")
				return
			}
		}
	}
}
