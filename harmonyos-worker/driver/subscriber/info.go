package subscriber

// SInfo 消息订阅信息
type SInfo struct {
	ExampleMsg MsgInfo `json:"example"`
	SystemMsg  MsgInfo `json:"system"`
}

// MsgInfo 消息信息
type MsgInfo struct {
	Topic string `json:"topic"`
	Queue string `json:"queue"`
}
