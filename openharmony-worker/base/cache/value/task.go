package value

// TaskValue 任务信息
type TaskValue struct {
	PID      string `json:"pid"`
	UserID   string `json:"userid"`
	DeviceID string `json:"deviceid"`
	ClientID string `json:"clientid"`
}
