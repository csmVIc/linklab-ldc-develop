package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func (sd *Driver) callback(m *nats.Msg) {
	var replyMsg msg.ReplyMsg
	defer func() {
		replyMsgByte, err := json.Marshal(replyMsg)
		if err != nil {
			sd.callbackErr = fmt.Errorf("natsconn.QueueSubscribe callback json.Marshal {%v}", err)
			log.Error(sd.callbackErr)
			return
		}
		err = m.Respond(replyMsgByte)
		if err != nil {
			sd.callbackErr = fmt.Errorf("natsconn.QueueSubscribe callback m.Respond {%v}", err)
			log.Error(sd.callbackErr)
			return
		}
	}()

	var tasks msg.DeviceBurnTasksMsg
	err := json.Unmarshal(m.Data, &tasks)
	if err != nil {
		sd.callbackErr = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	groupid, err := sd.createGroupID(tasks.UserID)
	if err != nil {
		sd.callbackErr = fmt.Errorf("create groupid error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	tasks.GroupID = groupid

	userTenantID, err := cache.Cdriver.GetUserTenantID(tasks.UserID)
	if err != nil {
		sd.callbackErr = fmt.Errorf("get user tenantid error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	tasks.TenantID = userTenantID

	tasksByte, err := json.Marshal(tasks)
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis hset error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	// 特殊设定
	if tasks.TenantID == 10001 {
		_, err = rdb.RPush(context.TODO(), "tasks:queue", string(tasksByte)).Result()
		if err != nil {
			sd.callbackErr = fmt.Errorf("redis rpush error {%v}", err)
			log.Error(sd.callbackErr)
			replyMsg.Code = -1
			replyMsg.Msg = err.Error()
			return
		}
	} else {
		_, err = rdb.LPush(context.TODO(), "tasks:queue", string(tasksByte)).Result()
		if err != nil {
			sd.callbackErr = fmt.Errorf("redis lpush error {%v}", err)
			log.Error(sd.callbackErr)
			replyMsg.Code = -1
			replyMsg.Msg = err.Error()
			return
		}
	}

	replyMsg.Code = 0
	replyMsg.Msg = "successfully entered the waiting queue"
	replyMsg.Data = map[string]string{
		"groupid": groupid,
	}

	// 记录进入等待队列日志
	tags := map[string]string{
		"userid": tasks.UserID,
	}
	fields := map[string]interface{}{
		"groupid":   groupid,
		"entertime": time.Now().UnixNano(),
		"podname":   os.Getenv("POD_NAME"),
		"nodename":  os.Getenv("NODE_NAME"),
	}
	err = logger.Ldriver.WriteLog("entertaskswait", tags, fields)
	if err != nil {
		sd.callbackErr = fmt.Errorf("write log err {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
}

func (sd *Driver) groupCallback(m *nats.Msg) {
	var replyMsg msg.ReplyMsg
	defer func() {
		replyMsgByte, err := json.Marshal(replyMsg)
		if err != nil {
			sd.callbackErr = fmt.Errorf("natsconn.QueueSubscribe callback json.Marshal {%v}", err)
			log.Error(sd.callbackErr)
			return
		}
		err = m.Respond(replyMsgByte)
		if err != nil {
			sd.callbackErr = fmt.Errorf("natsconn.QueueSubscribe callback m.Respond {%v}", err)
			log.Error(sd.callbackErr)
			return
		}
	}()

	var taskMsg msg.GroupBurnTaskMsg
	err := json.Unmarshal(m.Data, &taskMsg)
	if err != nil {
		sd.callbackErr = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	groupid, err := sd.createGroupID(taskMsg.UserID)
	if err != nil {
		sd.callbackErr = fmt.Errorf("create groupid error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	taskMsg.GroupID = groupid

	taskMsgByte, err := json.Marshal(taskMsg)
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis hset error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	_, err = rdb.LPush(context.TODO(), "tasks:group:queue", string(taskMsgByte)).Result()
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis lpush error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	replyMsg.Code = 0
	replyMsg.Msg = "successfully entered the waiting queue"
	replyMsg.Data = map[string]string{
		"groupid": groupid,
	}

	// 记录进入等待队列日志
	tags := map[string]string{
		"userid": taskMsg.UserID,
	}
	fields := map[string]interface{}{
		"groupid":   groupid,
		"entertime": time.Now().UnixNano(),
		"podname":   os.Getenv("POD_NAME"),
		"nodename":  os.Getenv("NODE_NAME"),
	}
	err = logger.Ldriver.WriteLog("entertaskswait", tags, fields)
	if err != nil {
		sd.callbackErr = fmt.Errorf("write log err {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
}
