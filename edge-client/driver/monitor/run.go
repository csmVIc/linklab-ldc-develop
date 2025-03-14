package monitor

import (
	"errors"
	"time"
)

// Run 监控运行
func (md *Driver) Run() error {

	// auth token 初始化
	// go api.ADriver.TokenUpdate()
	// tokentimeout := time.Now().Add(time.Second * time.Duration(md.info.Token.InitTimeOut))
	// for len(api.ADriver.GetToken()) < 1 {
	// 	if time.Now().After(tokentimeout) {
	// 		err := fmt.Errorf("auth token init wait timeout {%vs} error", md.info.Token.InitTimeOut)
	// 		log.Error(err)
	// 		return err
	// 	}
	// }

	// 心跳保持
	go md.heartbeat()
	// 设备状态
	go md.device()
	// Pod状态
	go md.pod()
	// 节点资源
	go md.resource()

	for md.endrun == false {
		time.Sleep(time.Second)
	}

	return errors.New("device end run")
}
