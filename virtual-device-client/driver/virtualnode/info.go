package virtualnode

import (
	"linklab/device-control-v2/base-library/parameter/msg"
	"sync"
	"time"
)

// BoardsInfo 开发板信息
type BoardsInfo struct {
	BurnDelay  int `json:"burndelay"`
	LogBytes   int `json:"logbytes"`
	VirtualNum int `json:"virtualnum"`
}

// DeviceLogInfo 设备日志信息
type DeviceLogInfo struct {
	LogChanSize     int `json:"logchansize"`
	LogTimeOutMill  int `json:"logtimeoutmill"`
	ReadSleepMill   int `json:"readsleepmill"`
	TaskTimeOutMill int `json:"tasktimeoutmill"`
}

// VInfo 设备信息
type VInfo struct {
	DeviceLog DeviceLogInfo         `json:"devicelog"`
	Boards    map[string]BoardsInfo `json:"boards"`
}

// DeviceBusyStatus 设备工作状态
type DeviceBusyStatus string

const (
	// Burning 烧写状态
	Burning DeviceBusyStatus = "Burning"
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
