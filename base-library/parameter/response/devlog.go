package response

// DeviceLog 运行日志
type DeviceLog struct {
	DeviceID string   `json:"deviceid"`
	ClientID string   `json:"clientid"`
	Logs     []string `json:"logs"`
}
