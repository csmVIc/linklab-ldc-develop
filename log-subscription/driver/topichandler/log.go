package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/problemsystem"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (td *Driver) logupload(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.Log.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var p msg.DeviceLog
	if err := json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		log.Errorf("json Unmarshal error {%v}", err)
		return
	}

	// 获取redis
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return
	}
	// 根据groupid和taskindex获取信息
	taskValueStr, err := rdb.HGet(context.TODO(), fmt.Sprintf("tasks:groupid:%s", p.GroupID), fmt.Sprintf("%v", p.TaskIndex)).Result()
	if err != nil {
		err = fmt.Errorf("redis hget {%v} {%v} error {%v}", fmt.Sprintf("tasks:groupid:%s", p.GroupID), p.TaskIndex, err)
		log.Error(err)
		return
	}
	taskValue := value.TaskValue{}
	if err = json.Unmarshal([]byte(taskValueStr), &taskValue); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return
	}

	userMsg := msg.UserMsg{
		Code: 0,
		Type: msg.TaskMsg,
		Data: msg.TaskData{
			GroupID:   p.GroupID,
			TaskIndex: p.TaskIndex,
			Type:      msg.TaskLogMsg,
			Msg:       p.Msg,
			Data: map[string]string{
				"ctimestamp": strconv.FormatInt(p.TimeStamp, 10),
				"dtimestamp": strconv.FormatInt(time.Now().UnixNano(), 10),
			},
		},
	}

	if err := messenger.Mdriver.PubMsg(fmt.Sprintf(td.info.Log.MsgTopic, taskValue.UserID), userMsg); err != nil {
		log.Errorf("messenger.Mdriver.PubMsg error {%v}", err)
		return
	}

	if len(taskValue.PID) > 0 && strings.HasPrefix(p.Msg, "end") {
		log.Errorf("p.Msg {%v} HasPrefix {end}, so skep", p.Msg)
		return
	}

	// 插入日志
	ctimestamp := time.Unix(0, p.TimeStamp)
	ctimestampstr := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", ctimestamp.Year(), ctimestamp.Month(), ctimestamp.Day(), ctimestamp.Hour(), ctimestamp.Minute(), ctimestamp.Second())

	elem := bson.M{
		"logs": fmt.Sprintf("%s:%s", ctimestampstr, p.Msg),
	}

	// 组日志
	gfilter := &table.GroupLogFilter{
		GroupID: p.GroupID,
	}
	if _, err := database.Mdriver.PushElemToArray("grouplog", gfilter, elem); err != nil {
		log.Errorf("database.Mdriver.PushElemToArray {grouplog} error {%v}", err)
		return
	}

	// 设备日志
	dfilter := &table.DeviceLogFilter{
		WaitingID: problemsystem.ComputeWaitingID(p.GroupID, p.TaskIndex),
	}
	if _, err := database.Mdriver.PushElemToArray("devicelog", dfilter, elem); err != nil {
		log.Errorf("database.Mdriver.PushElemToArray {devicelog} error {%v}", err)
		return
	}
}
