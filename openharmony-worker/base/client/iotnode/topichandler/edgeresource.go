package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) edgenoderesourcerefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt edgenode resource refuse {%v}", p.Msg)
	log.Error(err)
}

// PubEdgeNodeResource 发布边缘资源
func (td *Driver) PubEdgeNodeResource(resource *msg.EdgeNodeResourceList) error {

	if resource == nil || len(resource.EdgeNodes) < 1 {
		return nil
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["edgenoderesource"].Pub, 2, *resource)
	if err != nil {
		err = fmt.Errorf("mqtt client pub msg error {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}
