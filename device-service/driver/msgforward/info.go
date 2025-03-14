package msgforward

// SubInfo 订阅消息参数
type SubInfo struct {
	MsgTopic  string `json:"msgtopic"`
	MqttTopic string `json:"mqtttopic"`
	ChanSize  int    `json:"chansize"`
}

// MInfo 数据包转发的配置参数
type MInfo struct {
	ThreadMultiple int     `json:"threadmultiple"`
	Burn           SubInfo `json:"burn"`
	Write          SubInfo `json:"write"`
	PodApply       SubInfo `json:"podapply"`
}
