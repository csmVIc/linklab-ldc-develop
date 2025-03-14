package virtualnode

// GetDevices 获取当前设备列表
func (vd *Driver) GetDevices() map[string]bool {

	devices := make(map[string]bool)
	for dev := range vd.getdevices() {
		devices[dev] = true
	}

	return devices
}

// getdevices 获取当前设备列表
func (vd *Driver) getdevices() map[string]*DeviceStatus {

	tmpdevices := map[string]*DeviceStatus{}
	vd.devices.devicesMap.Range(func(key interface{}, value interface{}) bool {
		tmpdevices[key.(string)] = value.(*DeviceStatus)
		return true
	})

	return tmpdevices
}
