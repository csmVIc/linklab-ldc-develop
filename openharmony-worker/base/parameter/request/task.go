package request

// DeviceBurnTask 设备烧写任务
type DeviceBurnTask struct {
	BoardName string `form:"boardname" json:"boardname" binding:"required"`
	DeviceID  string `form:"deviceid" json:"deviceid" binding:"required"`
	RunTime   int    `form:"runtime" json:"runtime" binding:"required"`
	FileHash  string `form:"filehash" json:"filehash" binding:"required"`
	ClientID  string `form:"clientid" json:"clientid" binding:"required"`
	TaskIndex int    `form:"taskindex" json:"taskindex" binding:"required"`
}

// DeviceBurnTasks 设备烧写列表
type DeviceBurnTasks struct {
	PID   string           `form:"pid" json:"pid"`
	Tasks []DeviceBurnTask `form:"tasks" json:"tasks" binding:"required"`
}
