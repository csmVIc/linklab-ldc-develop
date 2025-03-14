package sub

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/tool"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func (sd *Driver) podApplyCallback(m *nats.Msg) {

	var replyMsg msg.ReplyMsg
	defer func() {
		replyMsgByte, err := json.Marshal(replyMsg)
		if err != nil {
			sd.callbackErr = fmt.Errorf("json.Marshal {%v}", err)
			log.Error(sd.callbackErr)
			return
		}
		err = m.Respond(replyMsgByte)
		if err != nil {
			sd.callbackErr = fmt.Errorf("m.Respond {%v}", err)
			log.Error(sd.callbackErr)
			return
		}
	}()

	p := msg.UserPodApply{}
	err := json.Unmarshal(m.Data, &p)
	if err != nil {
		sd.callbackErr = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	// // 创建任务
	// groupid, err := sd.createGroupID(p.UserID)
	// if err != nil {
	// 	sd.callbackErr = fmt.Errorf("sd.createGroupID error {%s}", err)
	// 	log.Error(sd.callbackErr)
	// 	replyMsg.Code = -1
	// 	replyMsg.Msg = err.Error()
	// 	return
	// }

	// 查询pod信息
	podname, nodeselector, err := sd.getPodInfo(p.YamlHash)
	if err != nil {
		err = fmt.Errorf("sd.getPodInfo error {%v}", err)
		log.Error(err)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	// 任务分配
	clientid, nodeaddselector, err := sd.allocate(p.UserID, podname, nodeselector)
	if err != nil {
		err = fmt.Errorf("edge allocate error {%v}", err)
		log.Error(err)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	// 获取地址
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	setupVal, err := rdb.HGet(context.TODO(), "edgenodes:setup", clientid).Result()
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis hget error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	setupInfo := value.EdgeNodeSetupInfo{}
	if err := json.Unmarshal([]byte(setupVal), &setupInfo); err != nil {
		sd.callbackErr = fmt.Errorf("json.Unmarshal error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	// 命名空间
	namesapce := tool.CreateMD5(p.UserID)

	// 获取客户端Token时添加详细日志
	log.Debugf("正在获取客户端 ID: %s 的令牌", clientid)
	// 获取客户端Token
	ctoken, err := auth.GetTokenFromID("clients", clientid)
    if err != nil {
        log.Debugf("获取令牌失败. 错误信息: %v", err)
        sd.callbackErr = fmt.Errorf("auth.GetTokenFromID error {%s}", err)
        log.Error(sd.callbackErr)
        replyMsg.Code = -1
        replyMsg.Msg = err.Error()
        return
    }
	/*
	Namespace:cce85e723d6886c97e13c3a7adae9841 
    YamlHash:15c0082d58767246080706108bd29be084af1dc919dfbe452dc017e1446c420c 
    UseEdgeRegistry:true 
    CreateIngress:false 
    NodeAddSelector:raspberrypi4bextend
	*/
    log.Debugf("成功获取令牌: %s", ctoken)

	// 在调用API之前添加Token验证日志
	log.Debugf("检查Redis中令牌状态键: clients:id:EdgeClient-0:token:%s", ctoken)

	// 部署服务
	podApplyReq := request.EdgeClientPodApply{
		Namespace:       namesapce,
		YamlHash:        p.YamlHash,
		UseEdgeRegistry: p.UseEdgeRegistry,
		CreateIngress:   p.CreateIngress,
		NodeAddSelector: nodeaddselector,
	}
	clientResp, err := sd.callEdgeClientPodApply(setupInfo.PodApplyURL, ctoken, &podApplyReq)
	if err != nil {
		sd.callbackErr = fmt.Errorf("%s", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}

	replyData := map[string]interface{}{
		"clientid": clientid,
	}

	if clientResp.Data != nil {
		respData := clientResp.Data.(map[string]interface{})
		replyData["ingressmap"] = respData["ingressmap"]
	} else {
		replyData["ingressmap"] = nil
	}

	// 回复报文
	replyMsg.Code = clientResp.Code
	replyMsg.Msg = clientResp.Msg
	replyMsg.Data = replyData

	return
}

func (ad *Driver) callEdgeClientPodApply(podapplyurl string, clienttoken string, podapplyreq *request.EdgeClientPodApply) (*response.Response, error) {

	client := resty.New()

    // 请求前的调试信息
    log.Debugf("准备调用API，请求地址: %s", podapplyurl)
    log.Debugf("使用的认证令牌长度: %d", len(clienttoken))
    log.Debugf("请求内容: %+v", podapplyreq)
	
	resp, err := client.R().
		SetHeader("Authorization", clienttoken).
		SetBody(*podapplyreq).Post(podapplyurl)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return nil, err
	}

    // 响应调试信息
    log.Debugf("响应状态码: %d", resp.StatusCode())
    log.Debugf("响应头信息: %+v", resp.Header())
    log.Debugf("响应内容: %s", string(resp.Body()))

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
