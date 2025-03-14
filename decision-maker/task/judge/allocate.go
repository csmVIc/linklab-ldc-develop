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

func (jd *Driver) taskAllocate(taskStr string, rdb *redis.ClusterClient) (*map[string][]msg.UserMsg, *map[string][]msg.ClientBurnMsg, error) {
	// json解析
	var task msg.DeviceBurnTasksMsg
	err := json.Unmarshal([]byte(taskStr), &task)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal {%v} error {%s}", taskStr, err)
		log.Error(err)
		return nil, nil, err
	}

	userTenantID, err := cache.Cdriver.GetUserTenantID(task.UserID)
	if err != nil {
		err = fmt.Errorf("cache.Cdriver.GetUserTenantID {%v} error {%s}", task.UserID, err)
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
	// 获取空闲设备
	idledevMap, err := jd.idledevices(userTenantID)
	if err != nil {
		err = fmt.Errorf("idledevices error {%s}", err)
		log.Error(err)
		return nil, nil, err
	}
	// log.Debugf("idledevMap %v", *idledevMap)
	// 约束条件判断, requiredDeviceMap统计该烧写任务需要哪些类型设备
	requiredNoSpecDeviceMap := make(map[string]int)
	requiredSpecDeviceMap := make(map[string]DeviceIndex)
	for _, elem := range task.Tasks {
		if len(elem.ClientID) > 0 && len(elem.DeviceID) > 0 {
			if _, isOk := (*idledevMap)[elem.BoardName]; isOk == false {
				err := fmt.Errorf("idle board {%v} is empty", elem.BoardName)
				log.Error(err)
				return nil, nil, err
			}

			// findindex := -1
			// for index, freedevice := range (*idledevMap)[elem.BoardName] {
			// 	if freedevice.ClientID == elem.ClientID && freedevice.DeviceID == elem.DeviceID {
			// 		requiredSpecDeviceMap[fmt.Sprintf("%s:%s", freedevice.ClientID, freedevice.DeviceID)] = freedevice
			// 		findindex = index
			// 		break
			// 	}
			// }

			freedeviceptr := (*idledevMap)[elem.BoardName].Front()
			for ; freedeviceptr != nil; freedeviceptr = freedeviceptr.Next() {
				freedevice := freedeviceptr.Value.(*DeviceIndex)
				if freedevice.ClientID == elem.ClientID && freedevice.DeviceID == elem.DeviceID {
					requiredSpecDeviceMap[fmt.Sprintf("%s:%s", freedevice.ClientID, freedevice.DeviceID)] = *freedevice
					break
				}
			}

			if freedeviceptr == nil {
				err := fmt.Errorf("can't find free device {%v} {%v}", elem.ClientID, elem.DeviceID)
				log.Error(err)
				return nil, nil, err
			}

			(*idledevMap)[elem.BoardName].Remove(freedeviceptr)

			// (*idledevMap)[elem.BoardName] = remove((*idledevMap)[elem.BoardName], findindex)
		} else {
			if _, isOk := requiredNoSpecDeviceMap[elem.BoardName]; isOk == false {
				requiredNoSpecDeviceMap[elem.BoardName] = 0
			}
			requiredNoSpecDeviceMap[elem.BoardName]++
		}
	}

	// 检测当前空闲设备是否满足，约束条件
	for boardname, number := range requiredNoSpecDeviceMap {
		if _, isOk := (*idledevMap)[boardname]; isOk == false {
			err := fmt.Errorf("idle board {%v} is empty", boardname)
			log.Error(err)
			return nil, nil, err
		}
		if (*idledevMap)[boardname].Len() < number {
			err := fmt.Errorf("idle board {%v} number {%v} < required number {%v}", boardname, (*idledevMap)[boardname].Len(), number)
			log.Error(err)
			return nil, nil, err
		}
	}

	// 判题系统日志
	glog := table.GroupLog{
		GroupID:    task.GroupID,
		UserID:     task.UserID,
		WaitingIDs: []string{},
		PID:        task.PID,
		Logs:       []string{},
	}

	// 修改所分配的每个任务和相应设备的状态
	userMsgMap := make(map[string][]msg.UserMsg)
	clientMsgMap := make(map[string][]msg.ClientBurnMsg)
	for _, elem := range task.Tasks {

		devUseStatus := value.DeviceUseStatus{
			GroupID:   task.GroupID,
			UserID:    task.UserID,
			TaskIndex: elem.TaskIndex,
			IsBurned:  false,
			RunTime:   elem.RunTime,
		}
		devUseStatusByte, err := json.Marshal(devUseStatus)
		if err != nil {
			err = fmt.Errorf("json.Marshal error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}

		devUse := DeviceIndex{}
		if len(elem.ClientID) > 0 && len(elem.DeviceID) > 0 {
			devUse = requiredSpecDeviceMap[fmt.Sprintf("%s:%s", elem.ClientID, elem.DeviceID)]
		} else {
			ptr := (*idledevMap)[elem.BoardName].Front()
			devUse = *(ptr.Value.(*DeviceIndex))
			(*idledevMap)[elem.BoardName].Remove(ptr)
		}
		_, err = rdb.HSet(context.TODO(), fmt.Sprintf("devices:use:%s", devUse.ClientID), devUse.DeviceID, string(devUseStatusByte)).Result()
		if err != nil {
			err = fmt.Errorf("rdb hset error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}

		taskValueByte, err := json.Marshal(value.TaskValue{
			PID:      task.PID,
			UserID:   task.UserID,
			ClientID: devUse.ClientID,
			DeviceID: devUse.DeviceID,
		})
		if err != nil {
			err = fmt.Errorf("json.Marshal error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}
		_, err = rdb.HSet(context.TODO(), fmt.Sprintf("tasks:groupid:%s", task.GroupID), fmt.Sprintf("%v", elem.TaskIndex), string(taskValueByte)).Result()
		if err != nil {
			err = fmt.Errorf("rdb hset error {%s}", err)
			log.Error(err)
			return nil, nil, err
		}

		if _, isOk := userMsgMap[task.UserID]; isOk == false {
			userMsgMap[task.UserID] = make([]msg.UserMsg, 0)
		}
		if _, isOk := clientMsgMap[devUse.ClientID]; isOk == false {
			clientMsgMap[devUse.ClientID] = make([]msg.ClientBurnMsg, 0)
		}
		userMsgMap[task.UserID] = append(userMsgMap[task.UserID], msg.UserMsg{
			Code: 0,
			Type: msg.TaskMsg,
			Data: msg.TaskData{
				GroupID:   task.GroupID,
				TaskIndex: elem.TaskIndex,
				Type:      msg.TaskAllocateMsg,
				Msg:       fmt.Sprint("the task is successfully assigned to the required device"),
				Data: map[string]string{
					"deviceid": devUse.DeviceID,
					"clientid": devUse.ClientID,
				},
			},
		})
		clientMsgMap[devUse.ClientID] = append(clientMsgMap[devUse.ClientID], msg.ClientBurnMsg{
			GroupID:   task.GroupID,
			DeviceID:  devUse.DeviceID,
			TaskIndex: elem.TaskIndex,
			FileHash:  elem.FileHash,
			RunTime:   elem.RunTime,
		})

		// 创建设备日志数据表
		waitingid := problemsystem.ComputeWaitingID(task.GroupID, elem.TaskIndex)
		dlog := table.DeviceLog{
			UserID:    task.UserID,
			ClientID:  devUse.ClientID,
			DevPort:   devUse.DeviceID,
			WaitingID: waitingid,
			GroupID:   task.GroupID,
			PID:       task.PID,
			Logs:      []string{},
			IsEnd:     false,
			StartDate: time.Now(),
			EndDate:   time.Now().Add(time.Duration(elem.RunTime) * time.Second),
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
		"userid": task.UserID,
	}
	fields := map[string]interface{}{
		"groupid":      task.GroupID,
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
