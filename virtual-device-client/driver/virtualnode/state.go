package virtualnode

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (vd *Driver) setburn(devport string) error {

	value, isOk := vd.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = Burning
	devstatus.BeginTime = time.Now()
	vd.devices.devicesMap.Store(devport, devstatus)

	return nil
}

func (vd *Driver) setrun(devport string) error {

	value, isOk := vd.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = Running
	devstatus.BeginTime = time.Now()
	vd.devices.devicesMap.Store(devport, devstatus)

	return nil
}

func (vd *Driver) setidle(devport string) error {

	value, isOk := vd.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devices not contain {%v}", devport)
		log.Error(err)
		return err
	}

	devstatus := value.(*DeviceStatus)
	devstatus.BusyStatus = IdleState
	devstatus.BeginTime = time.Now()
	vd.devices.devicesMap.Store(devport, devstatus)

	return nil
}

// GetDeviceStatus 获取设备状态
func (vd *Driver) GetDeviceStatus(devport string) (*DeviceStatus, error) {

	value, isOk := vd.devices.devicesMap.Load(devport)
	if isOk == false {
		err := fmt.Errorf("devport not exist {%v} error", devport)
		log.Error(err)
		return nil, err
	}

	devicestatus := value.(*DeviceStatus)
	return devicestatus, nil
}
