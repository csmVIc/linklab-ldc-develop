package linuxhostnode

import (
	"linklab/device-control-v2/base-library/parameter/msg"
	"sync"
	"time"
)

// BoardsInfo 开发板信息
type BoardsInfo struct {
	VirtualNum int    `json:"virtualnum"`
	ExecCmd    string `json:"execcmd"`
}

// DeviceLogInfo 设备日志信息
type DeviceLogInfo struct {
	LogChanSize    int `json:"logchansize"`
	LogTimeOutMill int `json:"logtimeoutmill"`
}

// LInfo 设备信息
type LInfo struct {
	Workspace string                `json:"workspace"`
	DeviceLog DeviceLogInfo         `json:"devicelog"`
	Boards    map[string]BoardsInfo `json:"boards"`
}

// DeviceBusyStatus 设备工作状态
type DeviceBusyStatus string

const (
	// Running 运行状态
	Running DeviceBusyStatus = "Running"
	// IdleState 空闲状态
	IdleState DeviceBusyStatus = "IdleState"
	// LogError 日志错误
	LogError DeviceBusyStatus = "LogError"
)

// DeviceStatus 设备占用信息
type DeviceStatus struct {
	BurnInfo   *msg.ClientBurnMsg
	BeginTime  time.Time
	BusyStatus DeviceBusyStatus
	LogChan    chan *LogMsg
	Lock       sync.Mutex
}

// TaskRunEnd 任务运行结束
type TaskRunEnd struct {
	BurnInfo *msg.ClientBurnMsg
	TimeOut  int
	EndTime  time.Time
}

// Devices 设备列表
type Devices struct {
	devicesMap  sync.Map
	devicesLock sync.RWMutex
}

// LogMsg 日志信息
type LogMsg struct {
	Msg       string
	TimeStamp int64
}
