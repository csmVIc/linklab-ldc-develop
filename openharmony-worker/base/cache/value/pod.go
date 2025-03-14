package value

import "fmt"

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
	CreateTime int64           `json:"createtime"`
	Containers []ContainerInfo `json:"containers"`
}

func (pi *PodInfo) Hash() string {
	return fmt.Sprintf("%v:%v", pi.Namespace, pi.Name)
}

// ContainerResourceInfo 容器资源信息
type ContainerResourceInfo struct {
	Name   string `json:"name"`
	CpuUse int64  `json:"cpuuse"`
	MemUse int64  `json:"memuse"`
}

// PodResourceInfo Pod资源信息
type PodResourceInfo struct {
	Name       string                  `json:"name"`
	Namespace  string                  `json:"namespace"`
	Containers []ContainerResourceInfo `json:"containers"`
}

func (pi *PodResourceInfo) Hash() string {
	return fmt.Sprintf("%v:%v", pi.Namespace, pi.Name)
}
