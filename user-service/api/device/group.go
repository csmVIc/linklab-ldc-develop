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
	"go.mongodb.org/mongo-driver/bson"
)

// 创建设备组
func creategroup(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	var p request.BoardGroup
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if len(p.Boards) < 1 || len(p.Boards) > dinfo.Group.MaxBoardsLen {
		err := fmt.Errorf("parameter boards length{%v} < 1 || > %v error", len(p.Boards), dinfo.Group.MaxBoardsLen)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 循环检查
	boardNameMap := map[string]bool{}
	for _, boardname := range p.Boards {

		// 避免重复检查
		if boardNameMap[boardname] {
			continue
		}

		// 检查开发板名
		if err := database.Mdriver.DocExist("boards", table.BoardFilter{BoardName: boardname}); err != nil {
			err := fmt.Errorf("parameter boardname{%v} not exist error {%v}", boardname, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		boardNameMap[boardname] = true
	}

	// 更新数据
	_, err := database.Mdriver.ReplaceElem("boardgroups", table.BoardGroupFilter{
		Type: p.Type,
	}, table.BoardGroup{
		Type:   p.Type,
		Boards: p.Boards,
	})
	if err != nil {
		err := fmt.Errorf("database.Mdriver.ReplaceElem error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
}

// 显示所有已定义的设备组
func listdefinegroup(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	cursor, err := database.Mdriver.FindElem("boardgroups", bson.M{})
	if err != nil {
		err := fmt.Errorf("database.Mdriver.FindElem error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	groups := []table.BoardGroup{}
	if err := cursor.All(context.TODO(), &groups); err != nil {
		err := fmt.Errorf("cursor.All error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	res := response.BindGroupInfoList{
		Groups: []response.BindGroupInfo{},
	}
	for _, group := range groups {
		elem := response.BindGroupInfo{
			Type:   group.Type,
			Boards: group.Boards,
		}
		res.Groups = append(res.Groups, elem)
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: res})
}

// 关联设备组
func linkgroup(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	var p request.DeviceGroup
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if len(p.Devices) < 1 || len(p.Devices) > dinfo.Group.MaxBoardsLen {
		err := fmt.Errorf("parameter devices length{%v} < 1 || > %v error", len(p.Devices), dinfo.Group.MaxBoardsLen)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	BGFilter := table.BoardGroupFilter{Type: p.Type}
	boardgroup := table.BoardGroup{}
	if err := database.Mdriver.FindOneElem("boardgroups", BGFilter, &boardgroup); err != nil {
		err := fmt.Errorf("boardgroup{%v} find error {%v}", p.Type, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查数量
	if len(p.Devices) != len(boardgroup.Boards) {
		err := fmt.Errorf("devices length {%v} != {%v} error", len(p.Devices), len(boardgroup.Boards))
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

	// 依次检查
	for index, boardname := range boardgroup.Boards {

		// 检查是否为规定的开发板类型
		if strings.Index(p.Devices[index].DeviceID, boardname) < 0 {
			err = fmt.Errorf("devices[%v] {%v:%v} device type is not {%v} error", index, p.Devices[index].ClientID, p.Devices[index].DeviceID, boardname)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 检查设备存在
		if _, err := rdb.HGet(context.TODO(), fmt.Sprintf("devices:active:%s", p.Devices[index].ClientID), p.Devices[index].DeviceID).Result(); err != nil {
			err = fmt.Errorf("device {%v:%v} not exist", p.Devices[index].ClientID, p.Devices[index].DeviceID)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 创建绑定组ID
	deviceGroupID, err := cache.Cdriver.CreateDeviceBindGroupID(p.Type)
	if err != nil {
		err = fmt.Errorf("cache.Cdriver.CreateDeviceBindGroupID {%v} error {%v}", p.Type, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	valueBytes, err := json.Marshal(p)
	if err != nil {
		err = fmt.Errorf("json.Marshal error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 记录设备绑定组
	if _, err := rdb.HSet(context.TODO(), fmt.Sprintf("bind:group:type:%v", p.Type), deviceGroupID, string(valueBytes)).Result(); err != nil {
		err = fmt.Errorf("rdb.HSet error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: map[string]string{
		"id": deviceGroupID,
	}})
}

// 显示所有已关联设备组
func listlinkgroup(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
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

	bindGroupTypes, err := rdb.Keys(context.TODO(), "bind:group:type:*").Result()
	if err != nil {
		err = fmt.Errorf("rdb.Keys err {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	result := response.DeviceGroupInfoList{
		Groups: []response.DeviceGroupInfo{},
	}
	for _, bindGroupType := range bindGroupTypes {
		bindGroups, err := rdb.HGetAll(context.TODO(), bindGroupType).Result()
		if err != nil {
			err = fmt.Errorf("rdb.HGetAll err {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		for bindGroupID, valueStr := range bindGroups {
			bindGroup := value.DeviceBindGroup{}
			if err := json.Unmarshal([]byte(valueStr), &bindGroup); err != nil {
				err = fmt.Errorf("json.Unmarshal err {%v}", err)
				log.Error(err)
				c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
				return
			}

			elem := response.DeviceGroupInfo{
				ID:      bindGroupID,
				Devices: []response.DevInfoForGroup{},
			}

			for _, devinfo := range bindGroup.Devices {
				elem.Devices = append(elem.Devices, response.DevInfoForGroup{
					ClientID: devinfo.ClientID,
					DeviceID: devinfo.DeviceID,
				})
			}

			result.Groups = append(result.Groups, elem)
		}
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: result})
}

// 取消关联设备组
func unlinkgroup(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	var p request.DeviceGroupUnlink
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
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

	// 获取设备绑定组类型
	bindGroupType := p.ID[:strings.Index(p.ID, "-")]

	// 检查设备绑定组是否存在
	if isOk, err := rdb.HExists(context.TODO(), fmt.Sprintf("bind:group:type:%v", bindGroupType), p.ID).Result(); err != nil {
		err = fmt.Errorf("rdb.HExists err {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	} else if isOk {
		if _, err := rdb.HDel(context.TODO(), fmt.Sprintf("bind:group:type:%v", bindGroupType), p.ID).Result(); err != nil {
			err = fmt.Errorf("rdb.HDel err {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
}

// 申请设备组
func allocategroup(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	var p request.GroupBurnTask
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
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

	// 检查是否有可用的设备绑定组
	if isOk, err := rdb.HLen(context.TODO(), fmt.Sprintf("bind:group:type:%v", p.Type)).Result(); err != nil {
		err = fmt.Errorf("rdb.HLen err {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	} else if isOk < 1 {
		err = fmt.Errorf("no group bind type {%v} available err", p.Type)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	BGFilter := table.BoardGroupFilter{Type: p.Type}
	boardgroup := table.BoardGroup{}
	if err := database.Mdriver.FindOneElem("boardgroups", BGFilter, &boardgroup); err != nil {
		err := fmt.Errorf("boardgroup{%v} find error {%v}", p.Type, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 开发板长度检查
	if len(p.Devices) != len(boardgroup.Boards) {
		err := fmt.Errorf("devices length {%v} != {%v} error", len(p.Devices), len(boardgroup.Boards))
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 开发板类型，烧写文件检查
	for index, boardname := range boardgroup.Boards {

		// 不为规定的开发板类型
		if boardname != p.Devices[index].BoardName {
			err := fmt.Errorf("device[%v] type {%v} != {%v} error", index, p.Devices[index].BoardName, boardname)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 烧写文件检查
		ffilter := table.FileFilter{BoardName: p.Devices[index].BoardName, FileHash: p.Devices[index].FileHash}
		if err = database.Mdriver.DocExist("files", ffilter); err != nil {
			err = fmt.Errorf("device[%v] filehash {%v} not exist error {%v}", index, p.Devices[index].FileHash, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 检查运行时间
	if p.RunTime < dinfo.TaskRuntimeLimit.MinRuntime {
		p.RunTime = dinfo.TaskRuntimeLimit.MinRuntime
	}
	if p.RunTime > dinfo.TaskRuntimeLimit.MaxRuntime {
		p.RunTime = dinfo.TaskRuntimeLimit.MaxRuntime
	}

	// 组装消息
	taskMsg := msg.GroupBurnTaskMsg{
		GroupID: "",
		UserID:  userid,
		Type:    p.Type,
		RunTime: p.RunTime,
		Devices: []msg.GroupBurnDeviceInfoMsg{},
	}
	for _, elem := range p.Devices {
		taskMsg.Devices = append(taskMsg.Devices, msg.GroupBurnDeviceInfoMsg{
			BoardName: elem.BoardName,
			FileHash:  elem.FileHash,
		})
	}

	// 发送消息
	reply := msg.ReplyMsg{}
	if err := messenger.Mdriver.RequestMsg(dinfo.Group.Topic, taskMsg, time.Second*time.Duration(dinfo.Group.ReplyTimeOut), &reply); err != nil {
		err = fmt.Errorf("user {%v}, nats request err {%v}", userid, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 成功进入任务等待队列
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: reply.Msg, Data: reply.Data})
}
