package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/jlink-client/driver/iotnode"
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
		devChangeMap, err := iotnode.IDriver.GetDevicesChange()
		if err != nil {
			log.Errorf("get devices change error {%v}", err)
			return err
		}
		if len(devChangeMap) > 0 {
			// 如果设备列表变化,应该立即上传
			log.Debug("devices list change")
			for dev, flag := range devChangeMap {
				log.Debugf("device {%v} {%v}", dev, flag)
				if flag {
					// 如果有新增设备,应该启动新的日志监控协程
					go md.readlog(dev)
				}
			}
			// 更新设备列表
			if err := topichandler.TDriver.PubDeviceUpdate(devChangeMap); err != nil {
				log.Errorf("topichandler.TDriver.PubDeviceUpdate error {%v}", err)
				return err
			}
		}

		// 定时上传设备,即使设备列表没有变化,也上传设备列表
		if time.Now().After(devicelastupdatetime) {
			devices := iotnode.IDriver.GetDevices()
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
