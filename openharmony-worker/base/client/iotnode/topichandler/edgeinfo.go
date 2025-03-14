package topichandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) edgenodesetuprefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt edgenode setup refuse {%v}", p.Msg)
	log.Error(err)
}

// PubEdgeNodeSetup 发布边缘信息
func (td *Driver) PubEdgeNodeSetup(info *msg.EdgeNodeSetup) error {

	if info == nil {
		return errors.New("info is nil error")
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["edgenodesetup"].Pub, 2, *info)
	if err != nil {
		err = fmt.Errorf("mqtt client pub msg error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}
