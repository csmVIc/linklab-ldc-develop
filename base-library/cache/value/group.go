package value

// DevInfoForGroup 设备信息
type DevInfoForGroup struct {
	ClientID string `json:"clientid"`
	DeviceID string `json:"deviceid"`
}

// DeviceBindGroup 设备绑定组
type DeviceBindGroup struct {
	Type    string            `json:"type"`
	Devices []DevInfoForGroup `json:"devices"`
}
