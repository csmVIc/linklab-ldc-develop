package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/linuxhost-client/driver/linuxhostnode"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) device() error {

	defer func() {
		// 如果该函数退出,则将结束运行标志位设置为true
		md.endrun = true
		log.Error("device func set endrun")
	}()

	firstloop := true
	devicelastupdatetime := time.Now().Add(-time.Hour)

	// 定时上传设备列表
	for md.endrun == false {

		// 第一次启动日志读取
		if time.Now().After(devicelastupdatetime) {
			devices := linuxhostnode.LDriver.GetDevices()
			if firstloop {
				for dev := range devices {
					go md.readlog(dev)
				}
				firstloop = false
			}

			// 更新设备列表
			if err := topichandler.TDriver.PubDeviceUpdate(devices); err != nil {
				log.Errorf("topichandler.TDriver.PubDeviceUpdate error {%v}", err)
				return err
			}
			devicelastupdatetime = time.Now().Add(time.Duration(md.info.DeviceUpdate.TimeOut) * time.Second)
		}

		time.Sleep(time.Millisecond * time.Duration(md.info.DeviceUpdate.DetectIntervalMill))
	}
	return nil
}
