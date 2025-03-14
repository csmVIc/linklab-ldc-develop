package request

// DevLogDownload 运行日志下载
type DevLogDownload struct {
	GroupID  string `form:"groupid" json:"groupid" binding:"required"`
	ClientID string `form:"clientid" json:"clientid"`
	DeviceID string `form:"deviceid" json:"deviceid"`
}
