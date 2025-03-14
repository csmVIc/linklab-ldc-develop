package msg

// DeviceEnd 设备日志
type DeviceEnd struct {
	GroupID   string `form:"groupid" json:"groupid" binding:"required"`
	TaskIndex int    `form:"taskindex" json:"taskindex" binding:"required"`
	TimeOut   int    `form:"timeout" json:"timeout" binding:"required"`
	EndTime   int64  `form:"endtime" json:"endtime" binding:"required"`
}
