package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) edgenodeupdaterefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt edgenode update refuse {%v}", p.Msg)
	log.Error(err)
}

// PubEdgeNodeUpdate 发布设备
func (td *Driver) PubEdgeNodeUpdate(edgeNodes *msg.EdgeNodeStatusList) error {

	if edgeNodes == nil || (len(edgeNodes.Delete) < 1 && len(edgeNodes.HeartBeat) < 1) {
		return nil
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["edgenodeupdate"].Pub, 2, *edgeNodes)
	if err != nil {
		err = fmt.Errorf("mqtt client pub msg error {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}
