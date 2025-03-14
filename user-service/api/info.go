package api

import (
	"linklab/device-control-v2/user-service/api/client"
	"linklab/device-control-v2/user-service/api/device"
	"linklab/device-control-v2/user-service/api/edgenode"
	"linklab/device-control-v2/user-service/api/ws"
)

// ServerAddress 服务绑定的地址
type ServerAddress struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// AInfo 服务的配置参数
type AInfo struct {
	Address  ServerAddress  `json:"address"`
	Ws       ws.WInfo       `json:"ws"`
	Device   device.DInfo   `json:"device"`
	Client   client.CInfo   `json:"client"`
	EdgeNode edgenode.EInfo `json:"edgenode"`
}
