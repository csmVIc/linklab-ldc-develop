package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/problemsystem"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) deviceendrun(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.EndRun.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var errsig error
	defer func() {
		if errsig != nil {
			err := mqttclient.MDriver.PubMsg(fmt.Sprintf(td.info.EndRun.RefuseTopic, username, clientid), 0, msg.ReplyMsg{
				Code: -1,
				Msg:  errsig.Error(),
				Data: nil,
			})
			if err != nil {
				log.Errorf("mqttclient.MDriver.PubMsg error {%v}", err)
				return
			}
		}
	}()

	var p msg.DeviceEnd
	if err := json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		err = fmt.Errorf("json Unmarshal error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 获取redis
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 根据groupid和taskindex获取信息
	taskValueStr, err := rdb.HGet(context.TODO(), fmt.Sprintf("tasks:groupid:%s", p.GroupID), fmt.Sprintf("%v", p.TaskIndex)).Result()
	if err != nil {
		err = fmt.Errorf("redis hget {%v} {%v} error {%v}", fmt.Sprintf("tasks:groupid:%s", p.GroupID), p.TaskIndex, err)
		log.Error(err)
		// errsig = err
		return
	}
	taskValue := value.TaskValue{}
	if err = json.Unmarshal([]byte(taskValueStr), &taskValue); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 获取设备使用状态
	devUseStatusStr, err := rdb.HGet(context.TODO(), fmt.Sprintf("devices:use:%s", username), taskValue.DeviceID).Result()
	if err != nil {
		err = fmt.Errorf("redis hget error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}
	devUseStatus := value.DeviceUseStatus{}
	if err := json.Unmarshal([]byte(devUseStatusStr), &devUseStatus); err != nil {
		err = fmt.Errorf("json Unmarshal error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 删除设备使用状态
	_, err = rdb.HDel(context.TODO(), fmt.Sprintf("devices:use:%s", username), taskValue.DeviceID).Result()
	if err != nil {
		err = fmt.Errorf("redis hdel error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 删除任务状态
	_, err = rdb.HDel(context.TODO(), fmt.Sprintf("tasks:groupid:%s", p.GroupID), fmt.Sprintf("%v", p.TaskIndex)).Result()
	if err != nil {
		err = fmt.Errorf("redis hdel error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 设备结束运行结果发送至用户日志话题
	msgStr := "detect program print {end}, so release the device"
	if p.TimeOut == 1 {
		msgStr = fmt.Sprintf("detect program exceed time limit {%vs}, so release the device", devUseStatus.RunTime)
	} else if p.TimeOut == 2 {
		msgStr = fmt.Sprintf("device {%v} is abort during running, so release the device", taskValue.DeviceID)
	}
	userMsg := msg.UserMsg{
		Code: 0,
		Type: msg.TaskMsg,
		Data: msg.TaskData{
			GroupID:   p.GroupID,
			TaskIndex: p.TaskIndex,
			Type:      msg.TaskEndRunMsg,
			Msg:       msgStr,
		},
	}

	if err := messenger.Mdriver.PubMsg(fmt.Sprintf(td.info.EndRun.MsgTopic, taskValue.UserID), userMsg); err != nil {
		err = fmt.Errorf("messenger.Mdriver.PubMsg err {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 设置运行结束
	waitingid := problemsystem.ComputeWaitingID(p.GroupID, p.TaskIndex)
	filter := table.DeviceLogFilter{
		WaitingID: waitingid,
	}
	elem := table.DeviceLogSetEnd{
		IsEnd:   true,
		EndDate: time.Now(),
	}
	_, err = database.Mdriver.UpdateElem("devicelog", filter, elem)
	if err != nil {
		err = fmt.Errorf("database.Mdriver.UpdateElem {devicelog} error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 开始判题
	if len(taskValue.PID) > 0 {
		err := problemsystem.PDriver.Start(waitingid, taskValue.PID)
		if err != nil {
			err = fmt.Errorf("problemsystem.PDriver.Start error {%v}", err)
			log.Error(err)
		}
	}

	// 记录日志
	tags := map[string]string{
		"userid": taskValue.UserID,
	}
	fields := map[string]interface{}{
		"groupid":   p.GroupID,
		"clientid":  username,
		"deviceid":  taskValue.DeviceID,
		"taskindex": p.TaskIndex,
		"timeout":   p.TimeOut,
		"endtime":   p.EndTime,
		"podname":   os.Getenv("POD_NAME"),
		"nodename":  os.Getenv("NODE_NAME"),
	}
	err = logger.Ldriver.WriteLog("endrun", tags, fields)
	if err != nil {
		err = fmt.Errorf("write log err {%v}", err)
		log.Error(err)
		errsig = err
		return
	}
}
