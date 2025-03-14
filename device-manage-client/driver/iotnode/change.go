package iotnode

import (
	"linklab/device-control-v2/base-library/parameter/msg"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetDevicesChange 获取设备列表变化, 以及
func (id *Driver) GetDevicesChange() (map[string]bool, error) {

	// 加锁
	id.devices.devicesLock.Lock()
	defer func() {
		id.devices.devicesLock.Unlock()
	}()

	// 获取正常设备列表, 以及错误设备列表
	tmpdevices := id.getdevices()

	// 显示系统中接入的设备
	nowdevices, err := id.lsdevices()

	
	if err != nil {
		log.Errorf("lsdevices error {%v}", err)
		return map[string]bool{}, err
	}

	devChangeMap := make(map[string]bool)
	for dev := range nowdevices {
		if _, isOk := tmpdevices[dev]; isOk == false {
			log.Error("***********************add")
			devChangeMap[dev] = true
		}
	}

	for dev := range tmpdevices {
		if _, isOk := nowdevices[dev]; isOk == false {
			//给予10s的设备状态维持，排除由于XIUOS烧写（断电-上电）而导致设备状态改变的情况。
			startTime := time.Now()

			for {
					// 检查时间是否超过了10秒
					if time.Since(startTime) >= 10*time.Second {
						log.Error("##################delete")
						devChangeMap[dev] = false
						break
					}

					nowdevices_update, err := id.lsdevices()
					if err != nil {
						log.Errorf("lsdevices error {%v}", err)
						return map[string]bool{}, err
					}

					if _, isOk := nowdevices_update[dev]; isOk == false {
						log.Error("********device state check********")
						continue
					} else {
						break
					}
			}
		}
	}

	for dev, flag := range devChangeMap {
		if flag {
			serialport, err := id.openserial(dev)
			if err != nil {
				// 如果该设备连串口都无法打开的话,那就没有必要添加该设备
				log.Errorf("devport {%v} open error {%v}", dev, err)
				delete(devChangeMap, dev)
				continue
			}
			cmdchan := make(chan *msg.DeviceCmd, id.info.Chan.CmdSize)
			id.devices.devicesMap.Store(dev, &DeviceStatus{
				BurnInfo:   nil,
				BeginTime:  time.Now(),
				BusyStatus: IdleState,
				SerialPort: serialport,
				CmdChan:    &cmdchan,
				Lock:       sync.Mutex{},
			})
		} else {
			id.devices.devicesMap.Delete(dev)
		}
	}

	// 设备错误信号
	for noelem := false; noelem == false; {
		select {
		case deverr := <-id.devices.deviceErrorChan:
			log.Debugf("recv device {%v} error {%v} signal", deverr.DeviceID, deverr.Err)
			if _, loaded := id.devices.devicesMap.LoadAndDelete(deverr.DeviceID); loaded {
				devChangeMap[deverr.DeviceID] = false
				log.Debugf("delete device {%v} from device list", deverr.DeviceID)
			}
			continue
		default:
			noelem = true
		}
	}

	return devChangeMap, nil
}
