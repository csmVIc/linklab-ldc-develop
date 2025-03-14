package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) execerr(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.BurnResult.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var errsig error
	defer func() {
		if errsig != nil {
			err := mqttclient.MDriver.PubMsg(fmt.Sprintf(td.info.ExecErr.RefuseTopic, username, clientid), 0, msg.ReplyMsg{
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

	var p msg.ExecErr
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

	// 删除使用状态
	_, err = rdb.HDel(context.TODO(), fmt.Sprintf("devices:use:%s", username), taskValue.DeviceID).Result()
	if err != nil {
		err = fmt.Errorf("redis hdel error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}
	_, err = rdb.HDel(context.TODO(), fmt.Sprintf("tasks:groupid:%s", p.GroupID), fmt.Sprintf("%v", p.TaskIndex)).Result()
	if err != nil {
		err = fmt.Errorf("redis hdel error {%v}", err)
		log.Error(err)
		errsig = err
		return
	}

	// 执行结果发布至用户日志话题
	userMsg := msg.UserMsg{
		Code: -1,
		Type: msg.TaskMsg,
		Data: msg.TaskData{
			GroupID:   p.GroupID,
			TaskIndex: p.TaskIndex,
			Type:      msg.TaskExecErrMsg,
			Msg:       p.Msg,
		},
	}

	if err := messenger.Mdriver.PubMsg(fmt.Sprintf(td.info.ExecErr.MsgTopic, taskValue.UserID), userMsg); err != nil {
		err = fmt.Errorf("messenger.Mdriver.PubMsg err {%v}", err)
		log.Error(err)
		errsig = err
		return
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
		"msg":       p.Msg,
		"podname":   os.Getenv("POD_NAME"),
		"nodename":  os.Getenv("NODE_NAME"),
	}
	err = logger.Ldriver.WriteLog("execerr", tags, fields)
	if err != nil {
		err = fmt.Errorf("write log err {%v}", err)
		log.Error(err)
		errsig = err
		return
	}
}
