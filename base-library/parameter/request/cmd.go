package request

// DeviceCmd 设备命令
type DeviceCmd struct {
	Cmd      string `form:"cmd" json:"cmd" binding:"required"`
	DeviceID string `form:"deviceid" json:"deviceid" binding:"required"`
	ClientID string `form:"clientid" json:"clientid" binding:"required"`
}
