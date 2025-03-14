package device

// TaskNumLimitInfo 任务数量限制
type TaskNumLimitInfo struct {
	MaxTaskNum int `json:"maxtasknum"`
	MinTaskNum int `json:"mintasknum"`
}

// TaskRuntimeLimitInfo 任务运行时间限制
type TaskRuntimeLimitInfo struct {
	MinRuntime int `json:"minruntime"`
	MaxRuntime int `json:"maxruntime"`
}

// MsgInfo 消息队列的信息
type MsgInfo struct {
	Topic        string `json:"topic"`
	ReplyTimeOut int    `json:"replytimeout"`
}

// CmdInfo 控制命令的消息队列信息
type CmdInfo struct {
	Topic        string `json:"topic"`
	ReplyTimeOut int    `json:"replytimeout"`
}

// GroupInfo 设备组信息
type GroupInfo struct {
	MaxBoardsLen int    `json:"maxboardslen"`
	Topic        string `json:"topic"`
	ReplyTimeOut int    `json:"replytimeout"`
}

// RegistryInfo 边缘镜像仓库
type RegistryInfo struct{
	RegistryAddress string `json:"registryaddress"`
}

// DInfo 设备烧写信息
type DInfo struct {
	TaskNumLimit     TaskNumLimitInfo     `json:"tasknumlimit"`
	TaskRuntimeLimit TaskRuntimeLimitInfo `json:"taskruntimelimit"`
	Msg              MsgInfo              `json:"msg"`
	Cmd              CmdInfo              `json:"cmd"`
	Group            GroupInfo            `json:"group"`
	Registry         RegistryInfo         `json:"registry"`
}

// DeviceIndex 用于索引到某个具体的设备
type DeviceIndex struct {
	ClientID string `json:"clientid"`
	DeviceID string `json:"deviceid"`
	Idle     bool   `json:"idle"`
}

// GroupInfo_ 绑定组信息
type GroupInfo_ struct {
	ID      string        `json:"id"`
	Devices []DeviceIndex `json:"devices"`
}

func remove(devs []DeviceIndex, i int) []DeviceIndex {
	devs[len(devs)-1], devs[i] = devs[i], devs[len(devs)-1]
	return devs[:len(devs)-1]
}
