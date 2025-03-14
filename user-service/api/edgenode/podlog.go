package edgenode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/tool"
	"linklab/device-control-v2/base-library/wsconf"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func podlog(c *gin.Context) {

	// 解析用户名
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 命名空间
	namespace := tool.CreateMD5(userid)

	// 参数验证
	p := request.UserPodLog{}
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 获取redis
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 查询Pod状态是否存在
	podexist, err := rdb.HExists(context.TODO(), fmt.Sprintf("pods:active:%v", p.ClientID), fmt.Sprintf("%v:%v", namespace, p.Pod)).Result()
	if err != nil {
		err = fmt.Errorf("rdb hexists error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	if !podexist {
		err = fmt.Errorf("pod {%v} not exist error", p.Pod)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 查询客户端地址
	setupVal, err := rdb.HGet(context.TODO(), "edgenodes:setup", p.ClientID).Result()
	if err != nil {
		err = fmt.Errorf("rdb.HGet error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	setupInfo := value.EdgeNodeSetupInfo{}
	if err := json.Unmarshal([]byte(setupVal), &setupInfo); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 获取客户端token
	ctoken, err := auth.GetTokenFromID("clients", p.ClientID)
	if err != nil {
		err = fmt.Errorf("auth.GetTokenFromID error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 建立和客户端之间的长连接
	edgeReq := &request.EdgeClientPodLog{
		Namespace: namespace,
		Pod:       p.Pod,
		Container: p.Container,
	}
	edgeUrl := fmt.Sprintf("%v?%v", setupInfo.PodLogURL, edgeReq.QueryRaw())
	edgehandler, resp, err := websocket.DefaultDialer.Dial(edgeUrl, http.Header{"Authorization": []string{ctoken}})
	if err != nil {
		if resp != nil {
			respbody := make([]byte, resp.ContentLength)
			resp.Body.Read(respbody)
			err = fmt.Errorf("websocket.DefaultDialer.Dial error {%v} resp {%v}", err, string(respbody))
		} else {
			err = fmt.Errorf("websocket.DefaultDialer.Dial error {%v}", err)
		}
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	defer edgehandler.Close()

	// 建立回复
	protocols := websocket.Subprotocols(c.Request)
	var httpHeader http.Header = nil
	if len(protocols) > 0 {
		log.Debugf("websocket subprotocols {%v}", protocols)
		httpHeader = http.Header{
			"Sec-Websocket-Protocol": []string{protocols[0]},
		}
	}

	// 长连接
	userhander, err := wsconf.UpgraderGlobal.Upgrade(c.Writer, c.Request, httpHeader)
	if err != nil {
		err := fmt.Errorf("websocket upgrade error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	defer userhander.Close()

	// 消息转发
	podLogMsgForward(userhander, edgehandler)
}
