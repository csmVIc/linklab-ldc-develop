package device

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
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

func sendcmd(c *gin.Context, topic string, replytimeout int) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	p := request.DeviceCmd{}
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("user {%v}, get rdb err {%v}", userid, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 查找设备是否为活跃
	if _, err := rdb.HGet(context.TODO(), fmt.Sprintf("devices:active:%s", p.ClientID), p.DeviceID).Result(); err != nil {
		err = fmt.Errorf("device {%v:%v} not exist", p.ClientID, p.DeviceID)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 查找设备是否为该用户占用
	deviceusestatusstr, err := rdb.HGet(context.TODO(), fmt.Sprintf("devices:use:%s", p.ClientID), p.DeviceID).Result()
	if err != nil {
		err = fmt.Errorf("user {%v} can't send cmd to device {%v:%v}", userid, p.ClientID, p.DeviceID)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	deviceusestatus := value.DeviceUseStatus{}
	if err := json.Unmarshal([]byte(deviceusestatusstr), &deviceusestatus); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%s}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	if deviceusestatus.UserID != userid {
		err = fmt.Errorf("user {%v} can't send cmd to device {%v:%v}", userid, p.ClientID, p.DeviceID)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 查询开发板所属类型
	boardname := p.DeviceID[strings.LastIndex(p.DeviceID, "/")+1 : strings.LastIndex(p.DeviceID, "-")]
	boardfilter := &table.BoardFilter{BoardName: boardname}
	boardfind := &table.Board{}
	if err := database.Mdriver.FindOneElem("boards", boardfilter, boardfind); err != nil {
		err = fmt.Errorf("boardname {%v} not support error {%v}", boardname, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	boardtypefilter := &table.BoardTypeFilter{BoardType: boardfind.BoardType}
	boardtypefind := &table.BoardType{}
	if err := database.Mdriver.FindOneElem("boardtypes", boardtypefilter, boardtypefind); err != nil {
		err = fmt.Errorf("boardtype {%v} not support error {%v}", boardfind.BoardType, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查开发板所属类型是否允许发送命令
	if boardtypefind.AllowCmd == false {
		err = fmt.Errorf("boardtype {%v} can't send cmd error", boardfind.BoardType)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 数据包序列化
	deviceCmd := msg.DeviceCmd{
		Cmd:       p.Cmd,
		DeviceID:  p.DeviceID,
		GroupID:   deviceusestatus.GroupID,
		TaskIndex: deviceusestatus.TaskIndex,
	}
	reply := msg.ReplyMsg{}
	if err := messenger.Mdriver.RequestMsg(fmt.Sprintf(topic, p.ClientID), deviceCmd, time.Second*time.Duration(replytimeout), &reply); err != nil {
		err = fmt.Errorf("user {%v} request msg error {%v}", userid, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 成功回复
	c.SecureJSON(http.StatusOK, response.Response{Code: reply.Code, Msg: reply.Msg, Data: reply.Data})
}
