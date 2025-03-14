package topichandler

// LogInfo 日志参数
type LogInfo struct {
	MsgTopic    string `json:"msgtopic"`
	RefuseTopic string `json:"refusetopic"`
}

// TInfo mqtt消息处理信息
type TInfo struct {
	Log LogInfo `json:"log"`
}
