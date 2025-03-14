package sub

// TopicInfo 订阅话题
type TopicInfo struct {
	PodApply   string `json:"podapply"`
	ImageBuild string `json:"imagebuild"`
}

// SInfo 任务分配订阅信息
type SInfo struct {
	Topic                 TopicInfo `json:"topic"`
	MaxReconn             int       `json:"maxreconn"`
	ReconnInterval        int       `json:"reconninterval"`
	MaxCreateGroupIDRetry int       `json:"maxcreategroupidretry"`
}
