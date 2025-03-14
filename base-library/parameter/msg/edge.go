package msg

// EdgeNodeResource 边缘设备资源
type EdgeNodeResource struct {
	Name         string `json:"name"`
	CpuAll       int64  `json:"cpuall"`
	MemAll       int64  `json:"memall"`
	CpuUse       int64  `json:"cpuuse"`
	MemUse       int64  `json:"memuse"`
	NvidiaGpuAll int64  `json:"nvidiagpuall"`
}

// EdgeNodeStatus 边缘设备状态
type EdgeNodeStatus struct {
	Name         string            `json:"name"`
	Ready        string            `json:"ready"`
	Architecture string            `json:"architecture"`
	OSImage      string            `json:"osimage"`
	OS           string            `json:"os"`
	IpAddress    string            `json:"ipaddress"`
	Labels       map[string]string `json:"labels"`
}

// EdgeNodeStatusList 边缘设备状态列表
type EdgeNodeStatusList struct {
	HeartBeat []EdgeNodeStatus `json:"heartbeat"`
	Delete    []string         `json:"delete"`
}

// EdgeNodeResourceList 边缘设备资源列表
type EdgeNodeResourceList struct {
	EdgeNodes []EdgeNodeResource `json:"edgenodes"`
}

// EdgeNodeSetup 边缘配置
type EdgeNodeSetup struct {
	PodApplyURL   string `json:"podapplyurl"`
	PodLogURL     string `json:"podlogurl"`
	PodDeleteURL  string `json:"poddeleteurl"`
	PodExecURL    string `json:"podexecurl"`
	ImageBuildURL string `json:"imagebuildurl"`
}
