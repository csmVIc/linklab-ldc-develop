package topichandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) endrunrefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt device endrun refuse {%v}", p.Msg)
	log.Error(err)
}

// PubEndRun 发布烧写结果
func (td *Driver) PubEndRun(burninfo *msg.ClientBurnMsg, timeout int, endtime time.Time) error {
	if burninfo == nil {
		err := errors.New("burninfo nil error")
		log.Error(err)
		return err
	}

	parameter := msg.DeviceEnd{
		GroupID:   burninfo.GroupID,
		TaskIndex: burninfo.TaskIndex,
		TimeOut:   timeout,
		EndTime:   endtime.UnixNano(),
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["endrun"].Pub, 2, parameter)
	if err != nil {
		log.Errorf("mqtt client pub msg error {%v}", err)
		return err
	}

	return nil
}
