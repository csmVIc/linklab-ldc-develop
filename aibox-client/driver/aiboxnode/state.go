package aiboxnode

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (ad *Driver) setrun(devport string) error {

	value, isOk := ad.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = Running
	devstatus.BeginTime = time.Now()
	ad.devices.devicesMap.Store(devport, devstatus)

	return nil
}

func (ad *Driver) setidle(devport string) error {

	value, isOk := ad.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = IdleState
	devstatus.BeginTime = time.Now()
	ad.devices.devicesMap.Store(devport, devstatus)

	return nil
}

// GetDeviceStatus 获取设备状态
func (ad *Driver) GetDeviceStatus(devport string) (*DeviceStatus, error) {

	value, isOk := ad.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devport not exist {%v} error", devport)
		log.Error(err)
		return nil, err
	}

	devicestatus := value.(*DeviceStatus)
	return devicestatus, nil
}

// GetDeviceReadLogChan 获取设备日志读通道
func (ad *Driver) GetDeviceReadLogChan(devport string) (<-chan *LogMsg, error) {
	value, isOk := ad.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devport not exist {%v} error", devport)
		log.Error(err)
		return nil, err
	}

	devicestatus := value.(*DeviceStatus)
	return devicestatus.LogChan, nil
}

// getDeviceWriteLogChan 获取设备日志写通道
func (ad *Driver) getDeviceWriteLogChan(devport string) (chan<- *LogMsg, error) {
	value, isOk := ad.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devport not exist {%v} error", devport)
		log.Error(err)
		return nil, err
	}

	devicestatus := value.(*DeviceStatus)
	return devicestatus.LogChan, nil
}
