package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) execerrrefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt exec error refuse {%v}", p.Msg)
	log.Error(err)
}

// PubExecErr 发布错误执行
func (td *Driver) PubExecErr(groupid string, taskindex int, errmsg string) error {

	execerr := msg.ExecErr{
		GroupID:   groupid,
		TaskIndex: taskindex,
		Msg:       errmsg,
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["execerr"].Pub, 2, execerr)
	if err != nil {
		log.Errorf("mqtt client pub msg error {%v}", err)
		return err
	}

	return nil
}
