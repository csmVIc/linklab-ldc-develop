package topichandler

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) deviceburnsub(client mqtt.Client, mqttmsg mqtt.Message) {

	var err error = nil
	defer func() {
		if err != nil {
			*td.errchan <- err
		}
	}()

	var p msg.ClientBurnMsg
	if err = json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		log.Errorf("json Unmarshal error {%v}", err)
		return
	}

	select {
	case (*td.burnchan) <- &p:
		log.Debugf("client burn msg into burnchan {%v}", p)
		return
	case <-time.After(time.Second * time.Duration(td.info.DeviceBurn.ChanTimeOut)):
		err = fmt.Errorf("client burn msg into burnchan timeout error {%vs}", td.info.DeviceBurn.ChanTimeOut)
		log.Error(err)
		return
	}
}

func (td *Driver) burnresultrefuse(client mqtt.Client, mqttmsg mqtt.Message) {
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

	err = fmt.Errorf("mqtt device burnresult refuse {%v}", p.Msg)
	log.Error(err)
}

// PubBurnResult 发布烧写结果
func (td *Driver) PubBurnResult(burninfo *msg.ClientBurnMsg, burntime *BurnTime, success bool, msgstr string) error {

	successflag := -1
	if success == true {
		successflag = 1
	}

	parameter := msg.BurnResult{
		GroupID:               burninfo.GroupID,
		TaskIndex:             burninfo.TaskIndex,
		Success:               successflag,
		Msg:                   msgstr,
		BeginBurnTime:         burntime.BeginBurn.UnixNano(),
		EndBurnTime:           burntime.EndBurn.UnixNano(),
		BeginDownloadFileTime: burntime.BeginDownloadFile.UnixNano(),
		EndDownloadFileTime:   burntime.EndDownloadFile.UnixNano(),
	}

	err := mqttclient.MDriver.PubMsg((*td.topicMap)["deviceburn"].Pub, 2, parameter)
	if err != nil {
		log.Errorf("mqtt client pub msg error {%v}", err)
		return err
	}

	return nil
}
