package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) podupdaterefuse(client mqtt.Client, mqttmsg mqtt.Message) {
	var err error = nil
	defer func() {
		if err != nil {
			*td.errchan <- err
		}
	}()

	p := msg.ReplyMsg{}
	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		err = fmt.Errorf("json Unmarshal error {%v}", err)
		log.Error(err)
		return
	}

	err = fmt.Errorf("mqtt pod update refuse {%v}", p.Msg)
	log.Error(err)
}

func (td *Driver) PubPodUpdate(pods *msg.PodStatusList) error {

	if pods == nil || (len(pods.Delete) < 1 && len(pods.HeartBeat) < 1) {
		return nil
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["podupdate"].Pub, 2, *pods)
	if err != nil {
		err = fmt.Errorf("mqtt client pub msg error {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}
