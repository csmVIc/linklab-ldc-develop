package messenger

// ClientInfo 消息队列连接信息
type ClientInfo struct {
	URL       string `json:"url"`
	ClusterID string `json:"clusterid"`
	NeedStan  bool   `json:"needstan"`
}

// MInfo 消息队列信息
type MInfo struct {
	Client ClientInfo `json:"client"`
}
