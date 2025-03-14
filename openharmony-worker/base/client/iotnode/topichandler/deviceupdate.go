package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) deviceupdaterefuse(client mqtt.Client, mqttmsg mqtt.Message) {

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

	err = fmt.Errorf("mqtt device update refuse {%v}", p.Msg)
	log.Error(err)
}

// PubDeviceUpdate 发布设备状态更新
func (td *Driver) PubDeviceUpdate(devChangeMap map[string]bool) error {

	if len(devChangeMap) < 1 {
		return nil
	}

	devices := msg.DeviceStatus{
		HeartBeat: []string{},
		Sub:       []string{},
	}

	// flag == true, 设备增加/设备心跳
	// flag == false, 设备删除
	for dev, flag := range devChangeMap {
		if flag == true {
			devices.HeartBeat = append(devices.HeartBeat, dev)
		} else {
			devices.Sub = append(devices.Sub, dev)
		}
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["deviceupdate"].Pub, 2, devices)
	if err != nil {
		log.Errorf("mqtt client pub msg error {%v}", err)
		return err
	}

	return nil
}
