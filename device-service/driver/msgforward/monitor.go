package msgforward

import (
	"errors"
	"linklab/device-control-v2/base-library/messenger"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// Monitor 启动监控接收协程
func (md *Driver) Monitor() error {

	for index := 0; index < runtime.NumCPU()*md.info.ThreadMultiple; index++ {
		log.Debugf("msgforward {%v} start up", index)

		// 烧写消息转发句柄
		go md.burnhandler(md.info.Burn.ChanSize, md.info.Burn.MsgTopic, md.info.Burn.MqttTopic)

		// 命令写入转发句柄
		go md.cmdhandler(md.info.Write.ChanSize, md.info.Write.MsgTopic, md.info.Write.MqttTopic)

		// Pod部署写入转发句柄
		// go md.podapplyhandler(md.info.PodApply.ChanSize, md.info.PodApply.MsgTopic, md.info.PodApply.MqttTopic)
	}

	for {
		if isOk := messenger.Mdriver.GetClosed(); isOk == true {
			return errors.New("messenger.Mdriver.GetClosed true error")
		}
		if md.exitsignal == true {
			return errors.New("sub handler exit error")
		}
		time.Sleep(time.Second)
	}
}
