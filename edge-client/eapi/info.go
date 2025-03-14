package eapi

// ServerAddress 服务绑定的地址
type ServerAddress struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// PodLogInfo Pod日志配置信息
type PodLogInfo struct {
	OutChanSize int `json:"outchansize"`
}

// EInfo 服务的配置参数
type EInfo struct {
	Address ServerAddress `json:"address"`
	PodLog  PodLogInfo    `json:"podlog"`
}
