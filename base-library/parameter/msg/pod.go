package msg

import "fmt"

// ContainerStatus 容器状态
type ContainerStatus struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartcount"`
	Image        string `json:"image"`
}

// PodStatus Pod状态
type PodStatus struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	NodeName   string            `json:"nodename"`
	Ready      string            `json:"ready"`
	CreateTime int64             `json:"createtime"`
	Containers []ContainerStatus `json:"containers"`
}

func (ps *PodStatus) Hash() string {
	return fmt.Sprintf("%v:%v", ps.Namespace, ps.Name)
}

// PodKey Pod键
type PodKey struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (pk *PodKey) Hash() string {
	return fmt.Sprintf("%v:%v", pk.Namespace, pk.Name)
}

// PodStatusList Pod状态列表
type PodStatusList struct {
	HeartBeat []PodStatus `json:"heartbeat"`
	Delete    []PodKey    `json:"delete"`
}

// // PodApply Pod部署
// type PodApply struct {
// 	GroupID   string `json:"groupid"`
// 	Namespace string `json:"namespace"`
// 	YamlHash  string `json:"yamlhash"`
// 	RunTime   int    `json:"runtime"`
// }

// // PodApplyResult Pod部署结果
// type PodApplyResult struct {
// 	GroupID        string `json:"groupid"`
// 	Success        bool   `json:"success"`
// 	Msg            string `json:"msg"`
// 	BeginApplyTime int64  `json:"beginapplytime"`
// 	EndApplyTime   int64  `json:"endapplytime"`
// }

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
	Containers []ContainerResource `json:"containers"`
}

// PodResourceList Pod资源列表
type PodResourceList struct {
	Pods []PodResource `json:"pods"`
}

// UserPodApply Pod部署
type UserPodApply struct {
	UserID          string `json:"userid"`
	YamlHash        string `json:"yamlhash"`
	UseEdgeRegistry bool   `json:"useedgeregistry"`
	CreateIngress   bool   `json:"createingress"`
}

// PodExecType Pod执行类型
type PodExecType string

const (
	NormalPodExec PodExecType = "NormalPodExec"
	ErrorPodExec  PodExecType = "ErrorPodExec"
)

// EdgeClientPodExec 边缘Pod执行
type EdgeClientPodExec struct {
	Type PodExecType `json:"type"`
	Msg  string      `json:"msg"`
}
