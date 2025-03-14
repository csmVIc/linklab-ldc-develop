package api

import "linklab/device-control-v2/file-cache/api/fcache"

// ServerAddress 服务绑定的地址
type ServerAddress struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// AInfo 服务的配置参数
type AInfo struct {
	Address   ServerAddress `json:"address"`
	FileCache fcache.FCInfo `json:"filecache"`
}
