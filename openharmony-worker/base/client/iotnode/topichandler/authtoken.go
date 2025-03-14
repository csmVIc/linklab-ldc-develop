package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) authtokensub(client mqtt.Client, mqttmsg mqtt.Message) {
	var err error = nil
	defer func() {
		if err != nil {
			*td.errchan <- err
		}
	}()

	p := msg.ReplyMsg{}
	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		log.Errorf("json Unmarshal error {%v}", err)
		return
	}

	if p.Code != 0 {
		log.Errorf("auth token sub code {%v} != 0 error {%v}", p.Code, p.Msg)
		return
	}

	tokenMap := p.Data.(map[string]interface{})
	if _, isOk := tokenMap["token"]; isOk == false {
		err = fmt.Errorf("auth token sub payload not contain {token} {%v}", p)
		log.Error(err)
		return
	}

	*td.tokenchan <- tokenMap["token"].(string)
}

func (td *Driver) authtokenrefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt device auth token refuse {%v}", p.Msg)
	log.Error(err)
}
