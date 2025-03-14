package iotnode

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/jlink-client/driver/iotnode/boardcmd"
	"sync"
)

// Driver 负责设备操作
type Driver struct {
	info        *IInfo
	devices     Devices
	boardCmdMap map[string]boardcmd.BoardCommand
}

// Driver 设备操作
var (
	IDriver *Driver
)

func (id *Driver) init() error {

	// nRF52840 初始化
	if _, isOk := id.info.Commands["nRF52840"]; isOk == false {
		return errors.New("init commands not contain nRF52840")
	}
	id.boardCmdMap["nRF52840"] = &boardcmd.NRF52840{}
	id.boardCmdMap["nRF52840"].SetEmptyProgram(id.info.Commands["nRF52840"].EmptyProgram)
	id.boardCmdMap["nRF52840"].SetScanCmd(id.info.Commands["nRF52840"].ScanCmd)

	return nil
}

// New 创建设备操作
func New(i *IInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}
	id := &Driver{info: i, devices: Devices{
		deviceErrorChan: make(chan *DeviceError, i.DeviceError.ChanSize),
		devicesMap:      sync.Map{},
		devicesLock:     sync.RWMutex{},
	}, boardCmdMap: make(map[string]boardcmd.BoardCommand)}
	if err := id.init(); err != nil {
		return nil, fmt.Errorf("init error {%v}", err)
	}
	return id, nil
}
