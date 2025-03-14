package topichandler

import "time"

// TopicInfo 话题信息
type TopicInfo struct {
	Refuse string `json:"refuse"`
	Pub    string `json:"pub"`
	Sub    string `json:"sub"`
}

// DeviceBurnInfo 烧写任务信息
type DeviceBurnInfo struct {
	ChanTimeOut int `json:"chantimeout"`
}

// PodApplyInfo Pod部署任务信息
// type PodApplyInfo struct {
// 	ChanTimeOut int `json:"chantimeout"`
// }

// TInfo 消息话题信息
type TInfo struct {
	Topics     map[string]TopicInfo `json:"topics"`
	DeviceBurn DeviceBurnInfo       `json:"deviceburn"`
	// PodApply   PodApplyInfo         `json:"podapply"`
}

// BurnTime 烧写时间信息
type BurnTime struct {
	BeginDownloadFile time.Time
	EndDownloadFile   time.Time
	BeginBurn         time.Time
	EndBurn           time.Time
}
