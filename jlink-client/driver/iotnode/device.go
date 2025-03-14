package iotnode

import log "github.com/sirupsen/logrus"

// GetSupportBoards 获取系统支持的板子类型
func (id *Driver) GetSupportBoards() []string {

	keys := make([]string, 0, len(id.boardCmdMap))
	for key := range id.boardCmdMap {
		keys = append(keys, key)
	}

	return keys
}

// GetDevices 获取当前设备列表
func (id *Driver) GetDevices() map[string]bool {

	devices := make(map[string]bool)
	for dev := range id.getdevices() {
		devices[dev] = true
	}

	return devices
}

// getdevices 获取当前设备列表
func (id *Driver) getdevices() map[string]*DeviceStatus {

	tmpdevices := map[string]*DeviceStatus{}
	id.devices.devicesMap.Range(func(key interface{}, value interface{}) bool {
		tmpdevices[key.(string)] = value.(*DeviceStatus)
		return true
	})

	return tmpdevices
}

// lsdevices 获取系统检测到的设备
func (id *Driver) lsdevices() (map[string]bool, error) {

	devices := map[string]bool{}
	for board := range id.boardCmdMap {
		devmap, err := id.boardCmdMap[board].Scan()
		if err != nil {
			log.Errorf("board {%v} scan error {%v}", board, err)
			return map[string]bool{}, err
		}

		for devname := range devmap {
			devices[devname] = true
		}
	}

	return devices, nil
}
