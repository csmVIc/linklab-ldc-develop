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
	TaskSize int `json:"tasksize"`
	CmdSize  int `json:"cmdsize"`
}

// BurnInfo 烧写信息
type BurnInfo struct {
	ThreadMultiple            int `json:"threadmultiple"`
	MaxFileDownloadRetry      int `json:"maxfiledownloadretry"`
	FileDownloadRetryInterval int `json:"filedownloadretryinterval"`
}

// TaskInfo 任务信息
type TaskInfo struct {
	ThreadMultiple int `json:"threadmultiple"`
}

// CmdWriteInfo 命令写入信息
type CmdWriteInfo struct {
	ThreadMultiple int `json:"threadmultiple"`
}

// MInfo 监控信息
type MInfo struct {
	API          api.AInfo        `json:"api"`
	HeartBeat    HeartBeatInfo    `json:"heartbeat"`
	DeviceUpdate DeviceUpdateInfo `json:"deviceupdate"`
	Chan         ChanInfo         `json:"chan"`
	Token        TokenInfo        `json:"token"`
	Burn         BurnInfo         `json:"burn"`
	Task         TaskInfo         `json:"task"`
	CmdWrite     CmdWriteInfo     `json:"cmdwrite"`
}
