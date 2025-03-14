package msg

// ClientBurnMsg 客户端烧写命令报文
type ClientBurnMsg struct {
	GroupID   string `json:"groupid"`
	DeviceID  string `json:"deviceid"`
	TaskIndex int    `json:"taskindex"`
	FileHash  string `json:"filehash"`
	RunTime   int    `json:"runtime"`
}

// DeviceStatus 客户端用于更新设备状态
type DeviceStatus struct {
	HeartBeat []string `form:"heartbeat" json:"heartbeat" binding:"required"`
	Sub       []string `form:"sub" json:"sub" binding:"required"`
}
