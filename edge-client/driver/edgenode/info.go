package edgenode

import (
	"sync"
)

// K8SClientInfo K8S客户端信息
type K8SClientInfo struct {
	InCluster  bool   `json:"incluster"`
	KubeConfig string `json:"kubeconfig"`
}

// K8SPodInfo Pod配置信息
type K8SPodInfo struct {
	SystemNamespaces map[string]bool `json:"systemnamespaces"`
}

// PodLogInfo Pod日志配置信息
type PodLogInfo struct {
	ChanTimeOut int `json:"chantimeout"`
}

// ImageBuildInfo 镜像构建信息
type ImageBuildInfo struct {
	RegistryAddress  string `json:"registryaddress"`
	BuildDownloadURL string `json:"builddownloadurl"`
}

// IngressInfo Ingress信息
type IngressInfo struct {
	Domain string `json:"domain"`
}

// EInfo 边缘信息
type EInfo struct {
	K8SClient  K8SClientInfo  `json:"k8sclient"`
	K8SPod     K8SPodInfo     `json:"k8spod"`
	PodLog     PodLogInfo     `json:"podlog"`
	ImageBuild ImageBuildInfo `json:"imagebuild"`
	Ingress    IngressInfo    `json:"ingress"`
}

// EdgeNodes 设备列表
type EdgeNodes struct {
	edgeNodesMap  sync.Map
	edgeNodesLock sync.RWMutex
}

// Pods Pod列表
type Pods struct {
	podsMap  sync.Map
	podsLock sync.RWMutex
}
