package response

// DeviceStatus 设备状态
type DeviceStatus struct {
	BoardName string `json:"boardname"`
	DeviceID  string `json:"deviceid"`
	Busy      bool   `json:"busy"`
	ClientID  string `json:"clientid"`
	Index     int    `json:"index"`
}

// DeviceList 设备列表
type DeviceList struct {
	Devices []DeviceStatus `json:"devices"`
}
