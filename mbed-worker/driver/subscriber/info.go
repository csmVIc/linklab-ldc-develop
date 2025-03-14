package subscriber

// SInfo 消息订阅信息
type SInfo struct {
	Topic string `json:"topic"`
	Queue string `json:"queue"`
}
