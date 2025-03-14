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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func poddelete(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 命名空间
	namespace := tool.CreateMD5(userid)

	// HTTP POST 参数验证
	p := request.UserPodDelete{}
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
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

	// 调用服务
	podDeleteReq := request.EdgeClientPodDelete{
		Namespace: namespace,
		Pod:       p.Pod,
	}
	clientResp, err := callEdgeClientPodDelete(setupInfo.PodDeleteURL, ctoken, &podDeleteReq)
	if err != nil {
		err = fmt.Errorf("pod delete error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: clientResp.Code, Msg: clientResp.Msg})
}

func callEdgeClientPodDelete(poddeleteurl string, clienttoken string, poddeletereq *request.EdgeClientPodDelete) (*response.Response, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", clienttoken).
		SetBody(*poddeletereq).Post(poddeleteurl)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return nil, err
	}

	result := response.Response{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json.Unmarshal error {%v}", resp.StatusCode())
		log.Error(err)
		return nil, err
	}

	return &result, nil
}
