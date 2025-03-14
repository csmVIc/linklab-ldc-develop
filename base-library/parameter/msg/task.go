package msg

// DeviceBurnTaskMsg 设备烧写任务消息
type DeviceBurnTaskMsg struct {
	BoardName string `json:"boardname"`
	DeviceID  string `json:"deviceid"`
	RunTime   int    `json:"runtime"`
	FileHash  string `json:"filehash"`
	ClientID  string `json:"clientid"`
	TaskIndex int    `json:"taskindex"`
}

// DeviceBurnTasksMsg 设备烧写任务列表消息
type DeviceBurnTasksMsg struct {
	GroupID  string              `json:"groupid"`
	UserID   string              `json:"userid"`
	TenantID int                 `json:"tenantid"`
	PID      string              `json:"pid"`
	Tasks    []DeviceBurnTaskMsg `json:"tasks"`
}
