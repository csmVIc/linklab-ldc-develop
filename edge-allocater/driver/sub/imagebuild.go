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

func (sd *Driver) imageBuildCallback(m *nats.Msg) {

	var replyMsg msg.ReplyMsg
	defer func() {
		log.Debugf("准备发送回复消息 - 状态码: {%d}, 消息: {%v}", replyMsg.Code, replyMsg.Msg)
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
    // 解析请求数据
    log.Debugf("开始解析镜像构建请求数据")
	p := msg.UserImageBuild{}
	err := json.Unmarshal(m.Data, &p)
	if err != nil {
		sd.callbackErr = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
    log.Debugf("请求数据解析成功 - 用户ID: {%v}, 镜像名称: {%v}", p.UserID, p.ImageName)

	// // 创建任务
	// groupid, err := sd.createGroupID(p.UserID)
	// if err != nil {
	// 	sd.callbackErr = fmt.Errorf("sd.createGroupID error {%s}", err)
	// 	log.Error(sd.callbackErr)
	// 	replyMsg.Code = -1
	// 	replyMsg.Msg = err.Error()
	// 	return
	// }

	// 任务分配
	log.Debugf("分配任务 - 用户ID: {%v}, 节点选择器: {%v}", p.UserID, p.NodeSelector)
	clientid, _, err := sd.allocate(p.UserID, "edge-build", p.NodeSelector)
	if err != nil {
		err = fmt.Errorf("edge allocate error {%v}", err)
		log.Error(err)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
    log.Debugf("任务分配成功 - 客户端ID: {%v}", clientid)

	// 获取地址，获取redis连接
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		sd.callbackErr = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	// 获取边缘节点设置信息
	log.Debugf("获取边缘节点设置信息 - 客户端ID: {%v}", clientid)
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
	log.Debugf("边缘节点设置信息解析成功啦 - 镜像构建URL: {%v}", setupInfo.ImageBuildURL)
	log.Debugf("边缘节点设置信息解析成功 - 详细信息: %+v", setupInfo)

	// 命名空间
	namesapce := tool.CreateMD5(p.UserID)
    log.Debugf("生成命名空间: {%v}", namesapce)

	// 获取客户端Token
	log.Debugf("开始获取客户端Token - client ID：{%v}",clientid)
	ctoken, err := auth.GetTokenFromID("clients", clientid)
	if err != nil {
		sd.callbackErr = fmt.Errorf("auth.GetTokenFromID error {%s}", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	log.Debugf("客户端Token - ctoken: {%v}", ctoken)
	// 镜像打包
	imageBuildReq := request.EdgeClientImageBuild{
		Namespace:    namesapce,
		FileHash:     p.FileHash,
		ImageName:    p.ImageName,
		NodeSelector: p.NodeSelector,
	}
	log.Debugf("镜像构建请求参数 - 命名空间: {%v}, 文件Hash: {%v}, 镜像名称: {%v}", 
		imageBuildReq.Namespace, 
		imageBuildReq.FileHash, 
		imageBuildReq.ImageName)
	clientResp, err := sd.callEdgeClientImageBuild(setupInfo.ImageBuildURL, ctoken, &imageBuildReq)
	if err != nil {
		sd.callbackErr = fmt.Errorf("%s", err)
		log.Error(sd.callbackErr)
		replyMsg.Code = -1
		replyMsg.Msg = err.Error()
		return
	}
	log.Debugf("镜像构建请求发送成功 - 响应码: {%v}, 响应消息: {%v}", 
		clientResp.Code, 
		clientResp.Msg)

	// 回复报文
	replyMsg.Code = clientResp.Code
	replyMsg.Msg = clientResp.Msg
	replyMsg.Data = map[string]string{
		"clientid": clientid,
	}
	log.Debugf("镜像构建回调处理完成")


	return
}

func (ad *Driver) callEdgeClientImageBuild(imagebuildurl string, clienttoken string, imagebuildreq *request.EdgeClientImageBuild) (*response.Response, error) {
    log.Debugf("开始调用边缘客户端构建镜像")
    log.Debugf("请求URL: {%v}", imagebuildurl)
    log.Debugf("请求参数 - 命名空间: {%v}, 镜像名称: {%v}", 
        imagebuildreq.Namespace, 
        imagebuildreq.ImageName)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", clienttoken).
		SetBody(*imagebuildreq).Post(imagebuildurl)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return nil, err
	}
    log.Debugf("HTTP请求完成 - 状态码: {%d}", resp.StatusCode())

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
    log.Debugf("边缘客户端响应 - 状态码: {%v}, 消息: {%v}", result.Code, result.Msg)
	return &result, nil
}
