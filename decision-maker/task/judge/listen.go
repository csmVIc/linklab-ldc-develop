package judge

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/decision-maker/task/pub"
	"time"

	log "github.com/sirupsen/logrus"
)

func (jd *Driver) listen() error {
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return err
	}
	for {

		task, err := rdb.BRPop(context.TODO(), 0, "tasks:queue").Result()
		if err != nil {
			err = fmt.Errorf("redis listen error {%s}", err)
			log.Error(err)
			return err
		}
		log.Debugf("get from {tasks:queue} task {%v}", task)

		userMsgMap, clientMsgMap, err := jd.taskAllocate(task[1], rdb)
		if err != nil {
			err = fmt.Errorf("task allocate erorr {%v}", err)
			log.Error(err)
			// 分配失败，重新加入到等待队列
			time.Sleep(time.Millisecond * time.Duration(jd.info.SleepMill))

			taskinfo := msg.DeviceBurnTasksMsg{}
			err := json.Unmarshal([]byte(task[1]), &taskinfo)
			if err != nil {
				err = fmt.Errorf("json unmarshal error {%s}", err)
				log.Error(err)
				return err
			}

			if taskinfo.TenantID == 10001 {
				_, err = rdb.RPush(context.TODO(), "tasks:tmp", task[1]).Result()
				if err != nil {
					err = fmt.Errorf("redis lpush error {%s}", err)
					log.Error(err)
					return err
				}
			} else {
				_, err = rdb.LPush(context.TODO(), "tasks:tmp", task[1]).Result()
				if err != nil {
					err = fmt.Errorf("redis lpush error {%s}", err)
					log.Error(err)
					return err
				}
			}

			// 检查是否需要交换
			if err = jd.taskQueueSwap(rdb); err != nil {
				err = fmt.Errorf("jd.taskQueueSwap( erorr {%v}", err)
				log.Error(err)
				return err
			}

			continue
		}

		// 检查是否需要交换
		if err = jd.taskQueueSwap(rdb); err != nil {
			err = fmt.Errorf("jd.taskQueueSwap( erorr {%v}", err)
			log.Error(err)
			return err
		}

		// 发送用户消息数据包
		err = pub.PDriver.PubUserMsg(userMsgMap)
		if err != nil {
			err = fmt.Errorf("pub.PDriver.PubUserMsg erorr {%v}", err)
			log.Error(err)
			return err
		}
		// 发送客户端烧写数据包
		err = pub.PDriver.PubClientBurnMsg(clientMsgMap)
		if err != nil {
			err = fmt.Errorf("pub.PDriver.PubClientBurnMsg erorr {%v}", err)
			log.Error(err)
			return err
		}
	}
}

func (jd *Driver) groupListen() error {
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return err
	}

	for {
		taskMsg, err := rdb.BRPop(context.TODO(), 0, "tasks:group:queue").Result()
		if err != nil {
			err = fmt.Errorf("redis listen error {%s}", err)
			log.Error(err)
			return err
		}
		log.Debugf("get from {tasks:group:queue} task {%v}", taskMsg)

		userMsgMap, clientMsgMap, err := jd.groupTaskAllocate(taskMsg[1], rdb)
		if err != nil {
			err = fmt.Errorf("task allocate erorr {%v}", err)
			log.Error(err)
			// 分配失败，重新加入到等待队列
			time.Sleep(time.Millisecond * time.Duration(jd.info.SleepMill))
			_, err = rdb.LPush(context.TODO(), "tasks:group:queue", taskMsg[1]).Result()
			if err != nil {
				err = fmt.Errorf("redis lpush error {%s}", err)
				log.Error(err)
				return err
			}
			continue
		}

		// 发送用户消息数据包
		err = pub.PDriver.PubUserMsg(userMsgMap)
		if err != nil {
			err = fmt.Errorf("pub.PDriver.PubUserMsg erorr {%v}", err)
			log.Error(err)
			return err
		}
		// 发送客户端烧写数据包
		err = pub.PDriver.PubClientBurnMsg(clientMsgMap)
		if err != nil {
			err = fmt.Errorf("pub.PDriver.PubClientBurnMsg erorr {%v}", err)
			log.Error(err)
			return err
		}
	}
}
