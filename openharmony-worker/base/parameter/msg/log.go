package msg

// DeviceLog 设备日志
type DeviceLog struct {
	GroupID   string `json:"groupid"`
	TaskIndex int    `json:"taskindex"`
	Msg       string `json:"msg"`
	TimeStamp int64  `json:"timestamp"`
}
