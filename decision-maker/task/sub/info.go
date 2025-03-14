package sub

// SInfo 烧写任务的订阅信息
type SInfo struct {
	Topic                 string `json:"topic"`
	GroupTopic            string `json:"grouptopic"`
	MaxReconn             int    `json:"maxreconn"`
	ReconnInterval        int    `json:"reconninterval"`
	MaxCreateGroupIDRetry int    `json:"maxcreategroupidretry"`
}
