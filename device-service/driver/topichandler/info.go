package topichandler

// LogInfo 日志参数
type LogInfo struct {
	MsgTopic    string `json:"msgtopic"`
	RefuseTopic string `json:"refusetopic"`
}

// EndRunInfo 结束运行参数
type EndRunInfo struct {
	MsgTopic    string `json:"msgtopic"`
	RefuseTopic string `json:"refusetopic"`
}

// BurnResultInfo 烧写结果参数
type BurnResultInfo struct {
	MsgTopic    string `json:"msgtopic"`
	RefuseTopic string `json:"refusetopic"`
}

// ExecErrInfo 执行错误参数
type ExecErrInfo struct {
	MsgTopic    string `json:"msgtopic"`
	RefuseTopic string `json:"refusetopic"`
}

// TInfo mqtt消息处理信息
type TInfo struct {
	Log        LogInfo        `json:"log"`
	EndRun     EndRunInfo     `json:"endrun"`
	BurnResult BurnResultInfo `json:"burnresult"`
	ExecErr    ExecErrInfo    `json:"execerr"`
}
