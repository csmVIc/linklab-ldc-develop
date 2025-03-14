package monitor

import "linklab/device-control-v2/base-library/client/iotnode/api"

// HeartBeatInfo 客户端心跳信息
type HeartBeatInfo struct {
	TimeOut int `json:"timeout"`
}

// DeviceUpdateInfo 设备更新信息
type DeviceUpdateInfo struct {
	TimeOut            int `json:"timeout"`
	DetectIntervalMill int `json:"detectintervalmill"`
}

// TokenInfo 秘钥信息
type TokenInfo struct {
	InitTimeOut int `json:"inittimeout"`
}

// ChanInfo 通道信息
type ChanInfo struct {
	BurnSize int `json:"burnsize"`
}

// ExecInfo 执行信息
type ExecInfo struct {
	ThreadMultiple            int `json:"threadmultiple"`
	MaxFileDownloadRetry      int `json:"maxfiledownloadretry"`
	FileDownloadRetryInterval int `json:"filedownloadretryinterval"`
}

// MInfo 监控信息
type MInfo struct {
	API          api.AInfo        `json:"api"`
	HeartBeat    HeartBeatInfo    `json:"heartbeat"`
	DeviceUpdate DeviceUpdateInfo `json:"deviceupdate"`
	Chan         ChanInfo         `json:"chan"`
	Token        TokenInfo        `json:"token"`
	Exec         ExecInfo         `json:"exec"`
}
