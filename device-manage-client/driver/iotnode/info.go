package iotnode

import (
	"linklab/device-control-v2/base-library/parameter/msg"
	"sync"
	"time"

	"github.com/albenik/go-serial/v2"
)

// CommandInfo 命令信息
type CommandInfo struct {
	Burn         string `json:"burn"`
	Reset        string `json:"reset"`
	BaudRate     int    `json:"baudrate"`
	NetworkCmd   string `json:"networkcmd"`
	NetworkScan  string `json:"networkscan"`
	EmptyProgram string `json:"emptyprogram"`
	WifiSSID     string `json:"wifissid"`
	WifiPassword string `json:"wifipassword"`
}

// DeviceErrorInfo 设备错误信息
type DeviceErrorInfo struct {
	ChanSize int `json:"chansize"`
	TimeOut  int `json:"timeout"`
}

// BurnInfo 烧写信息
type BurnInfo struct {
	MaxRetryTimes int `json:"maxretrytimes"`
}

// DeviceLogInfo 设备日志信息
type DeviceLogInfo struct {
	LogChanSize        int `json:"logchansize"`
	LogTimeOutMill     int `json:"logtimeoutmill"`
	LogSendTimeOutMill int `json:"logsendtimeoutmill"`
	ReadSleepMill      int `json:"readsleepmill"`
	TaskTimeOutMill    int `json:"tasktimeoutmill"`
}

// ChanInfo 通道信息
type ChanInfo struct {
	CmdSize int `json:"cmdsize"`
}

// IInfo 设备信息
type IInfo struct {
	Commands    map[string]CommandInfo `json:"commands"`
	Burn        BurnInfo               `json:"burn"`
	DeviceError DeviceErrorInfo        `json:"deviceerror"`
	DeviceLog   DeviceLogInfo          `json:"devicelog"`
	Chan        ChanInfo               `json:"chan"`
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
	SerialPort *serial.Port
	CmdChan    *chan *msg.DeviceCmd
	Lock       sync.Mutex
}

// DeviceError 设备错误
type DeviceError struct {
	DeviceID string
	Err      error
}

// TaskRunEnd 任务运行结束
type TaskRunEnd struct {
	BurnInfo *msg.ClientBurnMsg
	TimeOut  int
	EndTime  time.Time
}

// Devices 设备列表
type Devices struct {
	deviceErrorChan chan *DeviceError
	devicesMap      sync.Map
	devicesLock     sync.RWMutex
}

// LogMsg 日志信息
type LogMsg struct {
	Msg       string
	TimeStamp int64
}
