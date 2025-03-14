package value

// DeviceUseStatus 设备使用状况
type DeviceUseStatus struct {
	GroupID   string `json:"groupid"`
	UserID    string `json:"userid"`
	TaskIndex int    `json:"taskindex"`
	IsBurned  bool   `json:"isburned"`
	RunTime   int    `json:"runtime"`
}
