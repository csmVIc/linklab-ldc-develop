package response

// EdgeNodeStatus 边缘节点状态
type EdgeNodeStatus struct {
	Name         string `json:"name"`
	Ready        string `json:"ready"`
	Architecture string `json:"architecture"`
	OSImage      string `json:"osimage"`
	OS           string `json:"os"`
	IpAddress    string `json:"ipaddress"`
	ClientID     string `json:"clientid"`
}

// EdgeNodeList 边缘节点列表
type EdgeNodeList struct {
	EdgeNodes []EdgeNodeStatus `json:"edgenodes"`
}

// EdgeNodeResource 边缘设备资源
type EdgeNodeResource struct {
	Name         string `json:"name"`
	CpuAll       int64  `json:"cpuall"`
	MemAll       int64  `json:"memall"`
	CpuUse       int64  `json:"cpuuse"`
	MemUse       int64  `json:"memuse"`
	NvidiaGpuAll int64  `json:"nvidiagpuall"`
	ClientID     string `json:"clientid"`
}

// EdgeNodeResourceList 边缘设备资源列表
type EdgeNodeResourceList struct {
	EdgeNodes []EdgeNodeResource `json:"edgenodes"`
}
