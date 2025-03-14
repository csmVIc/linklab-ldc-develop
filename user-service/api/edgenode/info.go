package edgenode

// PodApplyInfo Pod部署信息
type PodApplyInfo struct {
	Topic        string `json:"topic"`
	ReplyTimeOut int    `json:"replytimeout"`
}

// ImageBuildInfo 镜像构建信息
type ImageBuildInfo struct {
	Topic        string `json:"topic"`
	ReplyTimeOut int    `json:"replytimeout"`
}

// PodLogInfo Pod日志信息
type PodLogInfo struct {
	MsgForwardChanSize int `json:"msgforwardchansize"`
}

// EInfo 边缘节点信息
type EInfo struct {
	PodApply   PodApplyInfo   `json:"podapply"`
	ImageBuild ImageBuildInfo `json:"imagebuild"`
	PodLog     PodLogInfo     `json:"podlog"`
}
