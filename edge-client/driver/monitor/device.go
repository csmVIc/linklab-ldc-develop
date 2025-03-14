package monitor

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) device() error {

	defer func() {
		// 如果该函数退出,则将结束运行标志位设置为true
		md.endrun = true
		log.Error("device func set endrun")
	}()

	devicelastupdatetime := time.Now()
	for md.endrun == false {
		// 检查设备变化
		edgeNodesChange, err := edgenode.EDriver.GetEdgeNodesChange()
		if err != nil {
			err = fmt.Errorf("edgenode.EDriver.GetEdgeNodesChange error {%v}", err)
			log.Error(err)
			return err
		}

		if len(edgeNodesChange.Delete) > 0 || len(edgeNodesChange.HeartBeat) > 0 {
			// 设备列表发生变化
			if err := topichandler.TDriver.PubEdgeNodeUpdate(edgeNodesChange); err != nil {
				err = fmt.Errorf("topichandler.TDriver.PubEdgeNodeUpdate error {%v}", err)
				log.Error(err)
				return err
			}
		}

		// 定时上传设备,即使设备列表没有变化,也上传设备列表
		if time.Now().After(devicelastupdatetime) {
			edgeNodes := edgenode.EDriver.GetEdgeNodes()
			// 更新设备列表
			if err := topichandler.TDriver.PubEdgeNodeUpdate(edgeNodes); err != nil {
				err = fmt.Errorf("topichandler.TDriver.PubEdgeNodeUpdate error {%v}", err)
				log.Error(err)
				return err
			}
			devicelastupdatetime = time.Now().Add(time.Duration(md.info.DeviceUpdate.TimeOut) * time.Second)
		}

		time.Sleep(time.Millisecond * time.Duration(md.info.DeviceUpdate.DetectIntervalMill))
	}
	return nil
}
