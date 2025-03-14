package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) clientdisconnect(client mqtt.Client, mqttmsg mqtt.Message) {

	clientMap := map[string]interface{}{}
	if err := json.Unmarshal(mqttmsg.Payload(), &clientMap); err != nil {
		log.Errorf("json.Unmarshal error {%v}", err)
		return
	}

	username := clientMap["username"].(string)
	clientid := clientMap["clientid"].(string)
	log.Debugf("mqtt client disconn username {%v} clientid {%v}", username, clientid)

	// 日志记录
	signoutsuccess := false
	defer func() {
		tags := map[string]string{
			"username": username,
		}
		fields := map[string]interface{}{
			"clientid": clientid,
			"success":  signoutsuccess,
			"podname":  os.Getenv("POD_NAME"),
			"nodename": os.Getenv("NODE_NAME"),
		}
		err := logger.Ldriver.WriteLog("clientdisconnect", tags, fields)
		if err != nil {
			err = fmt.Errorf("write log err {%v}", err)
			log.Error(err)
			return
		}
	}()

	// 检查客户端的登录状态
	tkey, err := auth.CheckClientIDAndGetToken(username, clientid)
	if err != nil {
		log.Errorf("auth.CheckClientIDAndGetToken {%v} {%v} error {%s}", username, clientid, err)
		return
	}

	// 查找缓存
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		log.Errorf("redis get rdb error {%s}", err)
		return
	}

	// 如果相同则删除该token
	_, err = rdb.Del(context.TODO(), tkey).Result()
	if err != nil {
		log.Errorf("{%v} redis del error {%v}", username, err)
		return
	}

	// 删除活跃设备
	if _, err := rdb.Del(context.TODO(), fmt.Sprintf("devices:active:%s", username)).Result(); err != nil {
		log.Errorf("rdb.Del {%v} error {%v}", fmt.Sprintf("devices:active:%s", username), err)
		return
	}

	// 处理所有未完成的任务
	devUseStatusMap, err := rdb.HGetAll(context.TODO(), fmt.Sprintf("devices:use:%s", username)).Result()
	if err != nil {
		log.Errorf("rdb.HGetAll {%v} error {%v}", fmt.Sprintf("devices:use:%s", username), err)
		return
	}

	for deviceid, devUseStatusStr := range devUseStatusMap {
		devUseStatus := value.DeviceUseStatus{}
		if err := json.Unmarshal([]byte(devUseStatusStr), &devUseStatus); err != nil {
			log.Errorf("json.Unmarshal error {%v}", err)
			return
		}
		taskData := msg.TaskData{
			GroupID:   devUseStatus.GroupID,
			TaskIndex: devUseStatus.TaskIndex,
			Type:      msg.TaskEndRunMsg,
			Msg:       fmt.Sprintf("client {%v} is abort during running, so release the device {%v}", username, deviceid),
		}
		if devUseStatus.IsBurned == false {
			taskData.Type = msg.TaskBurnMsg
			taskData.Msg = fmt.Sprintf("client {%v} is abort during burning, so release the device {%v}", username, deviceid)
		}
		userMsg := msg.UserMsg{
			Code: -1,
			Type: msg.TaskMsg,
			Data: taskData,
		}

		// 删除任务状态
		_, err = rdb.HDel(context.TODO(), fmt.Sprintf("tasks:groupid:%s", devUseStatus.GroupID), fmt.Sprintf("%v", devUseStatus.TaskIndex)).Result()
		if err != nil {
			log.Errorf("redis hdel error {%v}", err)
			return
		}

		// 将客户端错误信息汇报给用户
		if err := messenger.Mdriver.PubMsg(fmt.Sprintf(td.info.Disconnect.UserLogTopic, devUseStatus.UserID), userMsg); err != nil {
			log.Errorf("messenger.Mdriver.PubMsg err {%v}", err)
			return
		}
	}

	if len(devUseStatusMap) > 0 {
		if _, err := rdb.Del(context.TODO(), fmt.Sprintf("devices:use:%s", username)).Result(); err != nil {
			log.Errorf("rdb.Del {%v} error {%v}", fmt.Sprintf("devices:use:%s", username), err)
			return
		}
	}

	signoutsuccess = true
	log.Infof("client {%v} {%v} del token/device/task success", username, clientid)

	return
}
