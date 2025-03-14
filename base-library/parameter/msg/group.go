package msg

// GroupBurnDeviceInfoMsg 设备绑定组设备设备信息
type GroupBurnDeviceInfoMsg struct {
	BoardName string `json:"boardname"`
	FileHash  string `json:"filehash"`
}

// GroupBurnTaskMsg 设备绑定组烧写任务信息
type GroupBurnTaskMsg struct {
	GroupID string                   `json:"groupid"`
	UserID  string                   `json:"userid"`
	Type    string                   `json:"type"`
	RunTime int                      `json:"runtime"`
	Devices []GroupBurnDeviceInfoMsg `json:"devices"`
}
