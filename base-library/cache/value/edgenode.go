package value

// EdgeNodeInfo 边缘节点信息
type EdgeNodeInfo struct {
	Ready        string            `json:"ready"`
	Architecture string            `json:"architecture"`
	OSImage      string            `json:"osimage"`
	OS           string            `json:"os"`
	IpAddress    string            `json:"ipaddress"`
	Labels       map[string]string `json:"labels"`
}

// EdgeNodeResourceInfo 边缘节点资源信息
type EdgeNodeResourceInfo struct {
	CpuAll       int64 `json:"cpuall"`
	MemAll       int64 `json:"memall"`
	CpuUse       int64 `json:"cpuuse"`
	MemUse       int64 `json:"memuse"`
	NvidiaGpuAll int64 `json:"nvidiagpuall"`
}

// EdgeNodeSetupInfo 边缘节点配置信息
type EdgeNodeSetupInfo struct {
	PodApplyURL   string `json:"podapplyurl"`
	PodLogURL     string `json:"podlogurl"`
	PodDeleteURL  string `json:"poddeleteurl"`
	PodExecURL    string `json:"podexecurl"`
	ImageBuildURL string `json:"imagebuildurl"`
}
