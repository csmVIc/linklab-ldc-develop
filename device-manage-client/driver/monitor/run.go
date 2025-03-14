package monitor

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"time"

	log "github.com/sirupsen/logrus"
)

// Run 监控运行
func (md *Driver) Run() error {

	// auth token 初始化
	go api.ADriver.TokenUpdate()
	tokentimeout := time.Now().Add(time.Second * time.Duration(md.info.Token.InitTimeOut))
	for len(api.ADriver.GetToken()) < 1 {
		if time.Now().After(tokentimeout) {
			err := fmt.Errorf("auth token init wait timeout {%vs} error", md.info.Token.InitTimeOut)
			log.Error(err)
			return err
		}
	}

	// 心跳保持
	go md.heartbeat()
	// 任务结束
	go md.taskstartup()
	// 任务接收
	go md.burnstartup()
	// 命令接收
	go md.cmdwritestartup()
	// 设备状态
	go md.device()

	for md.endrun == false {
		time.Sleep(time.Second)
	}

	return errors.New("device end run")
}
