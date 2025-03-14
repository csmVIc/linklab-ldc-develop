package device

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"

	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func deviceburn(c *gin.Context, mintasknum int, maxtasknum int, minruntime int, maxruntime int, topic string, replytimeout int) {
	// 日志记录
	// entertime := time.Now()

	// 获取用户id
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	var p request.DeviceBurnTasks
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 任务数量检查
	if len(p.Tasks) < mintasknum || len(p.Tasks) > maxtasknum {
		err := fmt.Errorf("user {%v} task length {%v} not in range [%v:%v] error", userid, len(p.Tasks), mintasknum, maxtasknum)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 任务检查
	taskindexMap := make(map[int]bool)
	taskMsg := msg.DeviceBurnTasksMsg{
		GroupID: "",
		UserID:  userid,
		PID:     p.PID,
		Tasks:   make([]msg.DeviceBurnTaskMsg, 0),
	}
	for index := 0; index < len(p.Tasks); index++ {
		task := &p.Tasks[index]

		// 检查任务列表中是否有重复的索引
		if _, isOk := taskindexMap[task.TaskIndex]; isOk {
			err := fmt.Errorf("user {%v} task duplicate index {%v} error", userid, task.TaskIndex)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
		taskindexMap[task.TaskIndex] = true

		// 检查运行时间是否在范围内
		if task.RunTime < minruntime {
			task.RunTime = minruntime
		}
		if task.RunTime > maxruntime {
			task.RunTime = maxruntime
		}

		// 若指定设备，则clientid和deviceid均需要指定
		if (len(task.ClientID) < 1 && len(task.DeviceID) > 0) || (len(task.ClientID) > 0 && len(task.DeviceID) < 1) {
			err := fmt.Errorf("user {%v} task {%v}, if specify device, need clientid {%v} and deviceid {%v} not nil", userid, task.TaskIndex, task.ClientID, task.DeviceID)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 检查DeviceId是否属于BoardName
		if len(task.DeviceID) > 1 {

			splitBeginIndex := strings.LastIndex(task.DeviceID, "/")
			if splitBeginIndex == -1 {
				err := fmt.Errorf("user {%v} task {%v}, deviceid split {/} error", userid, task.TaskIndex)
				log.Error(err)
				c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
				return
			}

			splitEndIndex := strings.LastIndex(task.DeviceID, "-")
			if splitEndIndex == -1 {
				err := fmt.Errorf("user {%v} task {%v}, deviceid split {-} error", userid, task.TaskIndex)
				log.Error(err)
				c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
				return
			}

			if task.BoardName != task.DeviceID[splitBeginIndex+1:splitEndIndex] {
				task.BoardName = task.DeviceID[splitBeginIndex+1 : splitEndIndex]
			}
		}

		// 检查BoardName是否系统已经支持
		bfilter := table.BoardFilter{BoardName: task.BoardName}
		err := database.Mdriver.DocExist("boards", bfilter)
		if err != nil {
			err = fmt.Errorf("user {%v} task {%v}, boardname {%v} not support error {%v}", userid, task.TaskIndex, task.BoardName, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 检查filehash指定的文件是否存在
		ffilter := table.FileFilter{BoardName: task.BoardName, FileHash: task.FileHash}
		err = database.Mdriver.DocExist("files", ffilter)
		if err != nil {
			err = fmt.Errorf("user {%v} task {%v}, board {%v} filehash {%v} not exist {%v}", userid, task.TaskIndex, task.BoardName, task.FileHash, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 检查完成
		taskMsg.Tasks = append(taskMsg.Tasks, msg.DeviceBurnTaskMsg{
			BoardName: task.BoardName,
			DeviceID:  task.DeviceID,
			RunTime:   task.RunTime,
			FileHash:  task.FileHash,
			ClientID:  task.ClientID,
			TaskIndex: task.TaskIndex,
		})
	}

	reply := msg.ReplyMsg{}
	if err := messenger.Mdriver.RequestMsg(topic, taskMsg, time.Second*time.Duration(replytimeout), &reply); err != nil {
		err = fmt.Errorf("user {%v}, natsconn request err {%v}", userid, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 成功进入任务等待队列
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: reply.Msg, Data: reply.Data})
}
