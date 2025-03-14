package judge

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/problemsystem"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func (jd *Driver) groupTaskAllocate(taskStr string, rdb *redis.ClusterClient) (*map[string][]msg.UserMsg, *map[string][]msg.ClientBurnMsg, error) {

	// json解析
	var taskMsg msg.GroupBurnTaskMsg
	err := json.Unmarshal([]byte(taskStr), &taskMsg)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal {%v} error {%s}", taskStr, err)
		log.Error(err)
		return nil, nil, err
	}

	userTenantID, err := cache.Cdriver.GetUserTenantID(taskMsg.UserID)
	if err != nil {
		err = fmt.Errorf("cache.Cdriver.GetUserTenantID {%v} error {%s}", taskMsg.UserID, err)
		log.Error(err)
		return nil, nil, err
	}

	// 设备加锁
	devlocktoken, err := cache.Cdriver.Lock("devices:lock")
	if err != nil {
		err = fmt.Errorf("cache.Cdriver.Lock {%v} error {%s}", "devices:lock", err)
		log.Error(err)
		return nil, nil, err
	}
	defer func() {
		err = cache.Cdriver.UnLock("devices:lock", devlocktoken)
		if err != nil {
			err = fmt.Errorf("cache.Cdriver.UnLock {%v} error {%v}", "devices:lock", err)
			log.Error(err)
		}
	}()

	// 获取空闲设备绑定组
	groupInfo, err := jd.idlegroup(userTenantID, taskMsg.Type)
	if err != nil {
		err = fmt.Errorf("idlegroup error {%s}", err)
		log.Error(err)
		return nil, nil, err
	}

	// 判题系统日志
	glog := table.GroupLog{
		GroupID:    taskMsg.GroupID,
		UserID:     taskMsg.UserID,
		WaitingIDs: []string{},
		PID:        "",
		Logs:       []string{},
	}

	// 修改所分配的每个任务和相应设备的状态
	userMsgMap := make(map[string][]msg.UserMsg)
	clientMsgMap := make(map[string][]msg.ClientBurnMsg)
	for index, elem := range taskMsg.Devices {

		// 设置设备占用状态
		devUseStatus := value.DeviceUseStatus{
			GroupID:   taskMsg.GroupID,
			UserID:    taskMsg.UserID,
			TaskIndex: index,
			IsBurned:  false,
			RunTime:   taskMsg.RunTime,
		}
		devUseStatusByte, err := json.Marshal(devUseStatus)
		if err != nil {
			err = fmt.Errorf("json.Marshal error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}
		_, err = rdb.HSet(context.TODO(), fmt.Sprintf("devices:use:%s", groupInfo.Devices[index].ClientID), groupInfo.Devices[index].DeviceID, string(devUseStatusByte)).Result()
		if err != nil {
			err = fmt.Errorf("rdb hset error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}

		// 设置任务状态
		taskValueByte, err := json.Marshal(value.TaskValue{
			PID:      "",
			UserID:   taskMsg.UserID,
			ClientID: groupInfo.Devices[index].ClientID,
			DeviceID: groupInfo.Devices[index].DeviceID,
		})
		if err != nil {
			err = fmt.Errorf("json.Marshal error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}
		_, err = rdb.HSet(context.TODO(), fmt.Sprintf("tasks:groupid:%s", taskMsg.GroupID), fmt.Sprintf("%v", index), string(taskValueByte)).Result()
		if err != nil {
			err = fmt.Errorf("rdb hset error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}

		// 准备用户端、设备端报文
		if _, isOk := userMsgMap[taskMsg.UserID]; isOk == false {
			userMsgMap[taskMsg.UserID] = make([]msg.UserMsg, 0)
		}
		if _, isOk := clientMsgMap[groupInfo.Devices[index].ClientID]; isOk == false {
			clientMsgMap[groupInfo.Devices[index].ClientID] = make([]msg.ClientBurnMsg, 0)
		}
		userMsgMap[taskMsg.UserID] = append(userMsgMap[taskMsg.UserID], msg.UserMsg{
			Code: 0,
			Type: msg.TaskMsg,
			Data: msg.TaskData{
				GroupID:   taskMsg.GroupID,
				TaskIndex: index,
				Type:      msg.TaskAllocateMsg,
				Msg:       fmt.Sprint("the task is successfully assigned to the required device"),
				Data: map[string]string{
					"deviceid": groupInfo.Devices[index].DeviceID,
					"clientid": groupInfo.Devices[index].ClientID,
				},
			},
		})
		clientMsgMap[groupInfo.Devices[index].ClientID] = append(clientMsgMap[groupInfo.Devices[index].ClientID], msg.ClientBurnMsg{
			GroupID:   taskMsg.GroupID,
			DeviceID:  groupInfo.Devices[index].DeviceID,
			TaskIndex: index,
			FileHash:  elem.FileHash,
			RunTime:   taskMsg.RunTime,
		})

		// 创建设备日志数据表
		waitingid := problemsystem.ComputeWaitingID(taskMsg.GroupID, index)
		dlog := table.DeviceLog{
			UserID:    taskMsg.UserID,
			ClientID:  groupInfo.Devices[index].ClientID,
			DevPort:   groupInfo.Devices[index].DeviceID,
			WaitingID: waitingid,
			GroupID:   taskMsg.GroupID,
			PID:       "",
			Logs:      []string{},
			IsEnd:     false,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Duration(taskMsg.RunTime) * time.Second),
		}
		_, err = database.Mdriver.InsertElem("devicelog", dlog)
		if err != nil {
			err = fmt.Errorf("database.Mdriver.InsertElem {devicelog} err {%v}", err)
			log.Error(err)
			return nil, nil, err
		}

		glog.WaitingIDs = append(glog.WaitingIDs, waitingid)
	}

	// 创建判题系统日志
	_, err = database.Mdriver.InsertElem("grouplog", glog)
	if err != nil {
		err = fmt.Errorf("database.Mdriver.InsertElem {grouplog} err {%v}", err)
		log.Error(err)
		return nil, nil, err
	}

	// 记录分配日志
	tags := map[string]string{
		"userid": taskMsg.UserID,
	}
	fields := map[string]interface{}{
		"groupid":      taskMsg.GroupID,
		"allocatetime": time.Now().UnixNano(),
		"podname":      os.Getenv("POD_NAME"),
		"nodename":     os.Getenv("NODE_NAME"),
	}
	err = logger.Ldriver.WriteLog("taskallocate", tags, fields)
	if err != nil {
		err = fmt.Errorf("write log err {%v}", err)
		log.Error(err)
		return nil, nil, err
	}

	return &userMsgMap, &clientMsgMap, nil
}
