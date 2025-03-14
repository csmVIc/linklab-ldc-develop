package msg

// DeviceCmd 设备命令
type DeviceCmd struct {
	Cmd       string `json:"cmd"`
	DeviceID  string `json:"deviceid"`
	GroupID   string `json:"groupid"`
	TaskIndex int    `json:"taskindex"`
}
