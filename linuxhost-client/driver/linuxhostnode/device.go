package linuxhostnode

// GetDevices 获取当前设备列表
func (ld *Driver) GetDevices() map[string]bool {

	devices := make(map[string]bool)
	for dev := range ld.getdevices() {
		devices[dev] = true
	}

	return devices
}

// getdevices 获取当前设备列表
func (ld *Driver) getdevices() map[string]*DeviceStatus {

	tmpdevices := map[string]*DeviceStatus{}
	ld.devices.devicesMap.Range(func(key interface{}, value interface{}) bool {
		tmpdevices[key.(string)] = value.(*DeviceStatus)
		return true
	})

	return tmpdevices
}
