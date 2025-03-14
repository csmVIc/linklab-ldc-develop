package response

// DevInfoForGroup 设备信息
type DevInfoForGroup struct {
	ClientID string `json:"clientid"`
	DeviceID string `json:"deviceid"`
}

// DeviceGroupInfo 设备绑定组信息
type DeviceGroupInfo struct {
	ID      string            `json:"id"`
	Devices []DevInfoForGroup `json:"devices"`
}

// DeviceGroupInfoList 设备绑定组列表
type DeviceGroupInfoList struct {
	Groups []DeviceGroupInfo `json:"groups"`
}

// BindGroupInfo 设备绑定组信息
type BindGroupInfo struct {
	Type   string   `json:"type" bson:"type"`
	Boards []string `json:"boards" bson:"boards"`
}

// BindGroupInfoList 设备绑定组列表
type BindGroupInfoList struct {
	Groups []BindGroupInfo `json:"groups"`
}
