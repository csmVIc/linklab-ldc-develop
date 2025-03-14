package response

// ContainerInfo 边缘容器信息
type ContainerInfo struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartcount"`
	Image        string `json:"image"`
}

// PodInfo 边缘Pod信息
type PodInfo struct {
	Name       string          `json:"name"`
	Namespace  string          `json:"namespace"`
	NodeName   string          `json:"nodename"`
	Ready      string          `json:"ready"`
	ClientID   string          `json:"clientid"`
	CreateTime int64           `json:"createtime"`
	Containers []ContainerInfo `json:"containers"`
}

// PodList 边缘Pod列表
type PodList struct {
	Pods []PodInfo `json:"pods"`
}

// ContainerResource 容器资源
type ContainerResource struct {
	Name   string `json:"name"`
	CpuUse int64  `json:"cpuuse"`
	MemUse int64  `json:"memuse"`
}

// PodResource Pod资源
type PodResource struct {
	Name       string              `json:"name"`
	Namespace  string              `json:"namespace"`
	ClientID   string              `json:"clientid"`
	Containers []ContainerResource `json:"containers"`
}

// PodResourceList Pod资源列表
type PodResourceList struct {
	Pods []PodResource `json:"pods"`
}

// PodLogType Pod日志类型
type PodLogType string

const (
	NormalPodLog PodLogType = "NormalPodLog"
	ErrorPodLog  PodLogType = "ErrorPodLog"
)

// EdgeClientPodLog 边缘Pod日志
type EdgeClientPodLog struct {
	Type PodLogType `json:"type"`
	Msg  string     `json:"Msg"`
}

// EdgeClientPodApply Pod部署
type EdgeClientPodApply struct {
	IngressMap map[string]string `json:"ingressmap"`
}

// PodApply Pod部署
type PodApply struct {
	ClientID   string            `json:"clientid"`
	IngressMap map[string]string `json:"ingressmap"`
}
