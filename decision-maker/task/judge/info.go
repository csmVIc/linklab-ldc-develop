package judge

// JInfo 决策器的配置信息
type JInfo struct {
	SleepMill      int `json:"sleepmill"`
	MaxReconn      int `json:"maxreconn"`
	ReconnInterval int `json:"reconninterval"`
}

// DeviceIndex 用于索引到某个具体的设备
type DeviceIndex struct {
	ClientID string `json:"clientid"`
	DeviceID string `json:"deviceid"`
	Idle     bool   `json:"idle"`
}

// GroupInfo 绑定组信息
type GroupInfo struct {
	ID      string        `json:"id"`
	Devices []DeviceIndex `json:"devices"`
}

func remove(devs []DeviceIndex, i int) []DeviceIndex {
	devs[len(devs)-1], devs[i] = devs[i], devs[len(devs)-1]
	return devs[:len(devs)-1]
}
