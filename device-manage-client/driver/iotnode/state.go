package iotnode

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (id *Driver) setburn(devport string) error {

	value, isOk := id.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = Burning
	devstatus.BeginTime = time.Now()
	id.devices.devicesMap.Store(devport, devstatus)

	return nil
}

func (id *Driver) setrun(devport string) error {

	value, isOk := id.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = Running
	devstatus.BeginTime = time.Now()
	id.devices.devicesMap.Store(devport, devstatus)

	return nil
}

func (id *Driver) setidle(devport string) error {

	value, isOk := id.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = IdleState
	devstatus.BeginTime = time.Now()
	id.devices.devicesMap.Store(devport, devstatus)

	return nil
}

// GetDeviceStatus 获取设备状态
func (id *Driver) GetDeviceStatus(devport string) (*DeviceStatus, error) {

	value, isOk := id.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devport not exist {%v} error", devport)
		log.Error(err)
		return nil, err
	}

	devicestatus := value.(*DeviceStatus)
	return devicestatus, nil
}
