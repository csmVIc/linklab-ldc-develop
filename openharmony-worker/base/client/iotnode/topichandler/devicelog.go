package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) devicelogrefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt device log refuse {%v}", p.Msg)
	log.Error(err)
}

// PubDeviceLog 发布设备日志
func (td *Driver) PubDeviceLog(groupid string, taskindex int, logmsg string, timestamp int64) error {

	devlog := msg.DeviceLog{
		GroupID:   groupid,
		TaskIndex: taskindex,
		Msg:       logmsg,
		TimeStamp: timestamp,
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["devicelog"].Pub, 0, devlog)
	if err != nil {
		log.Errorf("mqtt client pub msg error {%v}", err)
		return err
	}

	return nil
}
