package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) heartbeatrefuse(client mqtt.Client, mqttmsg mqtt.Message) {

	var err error = nil
	defer func() {
		if err != nil {
			*td.errchan <- err
		}
	}()

	var p msg.ReplyMsg
	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		log.Errorf("json Unmarshal error {%v}", err)
		return
	}

	err = fmt.Errorf("mqtt heart beat refuse {%v}", p.Msg)
	log.Error(err)
}

// PubHeartBeat 发布客户端心跳消息
func (td *Driver) PubHeartBeat() error {

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["heartbeat"].Pub, 2, map[string]string{})
	if err != nil {
		log.Errorf("mqtt client pub msg error {%v}", err)
		return err
	}

	return nil
}
