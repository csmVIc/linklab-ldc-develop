package topichandler

import (
	"encoding/json"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) cmdwritesub(client mqtt.Client, mqttmsg mqtt.Message) {

	var err error = nil
	defer func() {
		if err != nil {
			*td.errchan <- err
		}
	}()

	p := msg.DeviceCmd{}
	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		log.Errorf("json Unmarshal error {%v}", err)
		return
	}

	*td.cmdchan <- &p
}
