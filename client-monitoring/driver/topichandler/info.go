package topichandler

// LoginInfo 登录参数
type LoginInfo struct {
	TTL              int    `json:"ttl"`
	TokenPubDelay    int    `json:"tokenpubdelay"`
	TokenPubTopic    string `json:"tokenpubtopic"`
	TokenRefuseTopic string `json:"tokenrefusetopic"`
}

// HeartBeatInfo 心跳参数
type HeartBeatInfo struct {
	TTL         int    `json:"ttl"`
	RefuseTopic string `json:"refusetopic"`
}

// DeviceInfo 设备参数
type DeviceInfo struct {
	TTL         int    `json:"ttl"`
	RefuseTopic string `json:"refusetopic"`
}

// EdgeNodeInfo 边缘节点参数
type EdgeNodeInfo struct {
	TTL         int    `json:"ttl"`
	RefuseTopic string `json:"refusetopic"`
}

// PodInfo 边缘Pod参数
type PodInfo struct {
	TTL         int    `json:"ttl"`
	RefuseTopic string `json:"refusetopic"`
}

// EdgeNodeResourceInfo 边缘节点资源参数
type EdgeNodeResourceInfo struct {
	TTL         int    `json:"ttl"`
	RefuseTopic string `json:"refusetopic"`
}

// EdgeNodeSetupInfo 边缘节点配置信息
type EdgeNodeSetupInfo struct {
	RefuseTopic string `json:"refusetopic"`
}

// PodResourceInfo 边缘Pod资源参数
type PodResourceInfo struct {
	TTL         int    `json:"ttl"`
	RefuseTopic string `json:"refusetopic"`
}

// DisconnectInfo 断开参数
type DisconnectInfo struct {
	UserLogTopic string `json:"userlogtopic"`
}

// TInfo mqtt消息处理信息
type TInfo struct {
	Login            LoginInfo            `json:"login"`
	HeartBeat        HeartBeatInfo        `json:"heartbeat"`
	Device           DeviceInfo           `json:"device"`
	Disconnect       DisconnectInfo       `json:"disconnect"`
	EdgeNode         EdgeNodeInfo         `json:"edgenode"`
	Pod              PodInfo              `json:"pod"`
	EdgeNodeResource EdgeNodeResourceInfo `json:"edgenoderesource"`
	PodResource      PodResourceInfo      `json:"podresource"`
	EdgeNodeSetup    EdgeNodeSetupInfo    `json:"edgenodesetup"`
}
