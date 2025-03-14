package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/api"
)

// HeartBeatInfo 客户端心跳信息
type HeartBeatInfo struct {
	TimeOut int `json:"timeout"`
}

// DeviceUpdateInfo 设备更新信息
type DeviceUpdateInfo struct {
	TimeOut            int `json:"timeout"`
	DetectIntervalMill int `json:"detectintervalmill"`
}

// PodUpdateInfo Pod更新信息
type PodUpdateInfo struct {
	TimeOut            int `json:"timeout"`
	DetectIntervalMill int `json:"detectintervalmill"`
}

// ResourseUpdateInfo 资源更新间隔
type ResourseUpdateInfo struct {
	Interval int `json:"interval"`
}

// TokenInfo 秘钥信息
type TokenInfo struct {
	InitTimeOut int `json:"inittimeout"`
}

// EdgeNodeSetupInfo 边缘节点配置
type EdgeNodeSetupInfo struct {
	Host string `json:"host"`
}

// MInfo 监控信息
type MInfo struct {
	API            api.AInfo          `json:"api"`
	HeartBeat      HeartBeatInfo      `json:"heartbeat"`
	DeviceUpdate   DeviceUpdateInfo   `json:"deviceupdate"`
	PodUpdate      PodUpdateInfo      `json:"podupdate"`
	ResourseUpdate ResourseUpdateInfo `json:"resourseupdate"`
	Token          TokenInfo          `json:"token"`
	EdgeNodeSetup  EdgeNodeSetupInfo  `json:"edgenodesetup"`
}
