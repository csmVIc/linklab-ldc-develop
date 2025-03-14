package edgenode

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
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func podapply(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	var p request.UserPodApply
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查filehash指定的文件是否存在
	filter := table.PodYamlFilter{
		FileHash: p.YamlHash,
	}
	if err := database.Mdriver.DocExist("podyaml", filter); err != nil {
		err = fmt.Errorf("podyaml {%v} not exist {%v}", p.YamlHash, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
	}

	// 消息报文
	applyMsg := msg.UserPodApply{
		UserID:          userid,
		YamlHash:        p.YamlHash,
		UseEdgeRegistry: p.UseEdgeRegistry,
		CreateIngress:   p.CreateIngress,
	}

	// 请求
	reply := msg.ReplyMsg{}
	if err := messenger.Mdriver.RequestMsg(einfo.PodApply.Topic, applyMsg, time.Second*time.Duration(einfo.PodApply.ReplyTimeOut), &reply); err != nil {
		err = fmt.Errorf("natsconn request err {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 回复
	c.SecureJSON(http.StatusOK, response.Response{Code: reply.Code, Msg: reply.Msg, Data: reply.Data})
}
