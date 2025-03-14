package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) heartbeat() error {

	defer func() {
		md.endrun = true
		log.Error("heartbeat func set endrun")
	}()

	for md.endrun == false {
		select {
		case err := <-md.errchan:
			// 接收到错误信息,断开mqtt连接
			log.Errorf("monitor heartbeat recv error {%v}", err)
			return err

		case <-time.After(time.Second * time.Duration(md.info.HeartBeat.TimeOut)):
			// 更新客户端心跳
			if err := topichandler.TDriver.PubHeartBeat(); err != nil {
				log.Errorf("topichandler.TDriver.PubHeartBeat error {%v}", err)
				return err
			}
		}
	}

	return nil
}
