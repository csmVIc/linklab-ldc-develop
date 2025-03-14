package monitor

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) TokenInit() error {
	// 1. 启动一个 goroutine 监听 token 更新
	go api.ADriver.TokenUpdate()
	// 2. 设置超时时间
	tokentimeout := time.Now().Add(time.Second * time.Duration(md.info.Token.InitTimeOut))
	// 3. 等待直到收到 token 或超时
	for len(api.ADriver.GetToken()) < 1 {
		if time.Now().After(tokentimeout) {
			err := fmt.Errorf("auth token init wait timeout {%vs} error", md.info.Token.InitTimeOut)
			log.Error(err)
			return err
		}
	}
	return nil
}

func (md *Driver) SetupInit() error {

	info := msg.EdgeNodeSetup{
		PodApplyURL:   fmt.Sprintf("http://%v/api/pod/apply", md.info.EdgeNodeSetup.Host),
		PodLogURL:     fmt.Sprintf("ws://%v/api/pod/log", md.info.EdgeNodeSetup.Host),
		PodDeleteURL:  fmt.Sprintf("http://%v/api/pod/delete", md.info.EdgeNodeSetup.Host),
		PodExecURL:    fmt.Sprintf("ws://%v/api/pod/exec", md.info.EdgeNodeSetup.Host),
		ImageBuildURL: fmt.Sprintf("http://%v/api/image/build", md.info.EdgeNodeSetup.Host),
	}

	if err := topichandler.TDriver.PubEdgeNodeSetup(&info); err != nil {
		err := fmt.Errorf("topichandler.TDriver.PubEdgeNodeSetup {%v} error", err)
		log.Error(err)
		return err
	}

	return nil
}
