package mqttclient

import mqtt "github.com/eclipse/paho.mqtt.golang"

// ClientInfo mqtt连接信息
type ClientInfo struct {
	IsCloud  bool   `json:"iscloud"`
	URL      string `json:"url"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

// MonitorInfo mqtt运行时信息
type MonitorInfo struct {
	MaxDisconnWait int `json:"maxdisconnwait"`
}

// PublishInfo 消息发布
type PublishInfo struct {
	TimeOut int `json:"timeout"`
}

// MInfo mqtt信息
type MInfo struct {
	Client  ClientInfo  `json:"client"`
	Monitor MonitorInfo `json:"monitor"`
	Publish PublishInfo `json:"publish"`
}

// SubTopic 订阅信息
type SubTopic struct {
	MsgHandler mqtt.MessageHandler
	Qos        byte
}
