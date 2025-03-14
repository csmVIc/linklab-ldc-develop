package request

// BoardGroup 开发板绑定组
type BoardGroup struct {
	Type   string   `form:"type" json:"type" binding:"required"`
	Boards []string `form:"boards" json:"boards" binding:"required"`
}

// DevInfoForGroup 设备信息
type DevInfoForGroup struct {
	ClientID string `form:"clientid" json:"clientid" binding:"required"`
	DeviceID string `form:"deviceid" json:"deviceid" binding:"required"`
}

// DeviceGroup 设备绑定组
type DeviceGroup struct {
	Type    string            `form:"type" json:"type" binding:"required"`
	Devices []DevInfoForGroup `form:"devices" json:"devices" binding:"required"`
}

// DeviceGroupUnlink 设备组取消关联
type DeviceGroupUnlink struct {
	ID string `json:"id"`
}

// GroupBurnDeviceInfo 组烧写任务设备信息
type GroupBurnDeviceInfo struct {
	BoardName string `form:"boardname" json:"boardname" binding:"required"`
	FileHash  string `form:"filehash" json:"filehash" binding:"required"`
}

// GroupBurnTask 组烧写任务
type GroupBurnTask struct {
	Type    string                `form:"type" json:"type" binding:"required"`
	RunTime int                   `form:"runtime" json:"runtime" binding:"required"`
	Devices []GroupBurnDeviceInfo `form:"devices" json:"devices" binding:"required"`
}
