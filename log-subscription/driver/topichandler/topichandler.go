package topichandler

import (
	"errors"
	"linklab/device-control-v2/base-library/mqttclient"

	log "github.com/sirupsen/logrus"
)

// Driver 负责mqtt消息处理
type Driver struct {
	info       *TInfo
	topicToSub *map[string]mqttclient.SubTopic
}

var (
	//TDriver mqtt消息处理全局实例
	TDriver *Driver
)

func (td *Driver) init() error {
	ts := make(map[string]mqttclient.SubTopic)
	ts["$share/device-service/clients/+/+/publish/log/upload"] = mqttclient.SubTopic{
		MsgHandler: td.logupload,
		Qos:        0,
	}
	td.topicToSub = &ts
	return nil
}

// New 创建mqtt消息处理实例
func New(i *TInfo) (*Driver, error) {
	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}
	td := &Driver{info: i, topicToSub: nil}
	if err := td.init(); err != nil {
		log.Errorf("topic handler init error {%v}", err)
		return nil, err
	}
	return td, nil
}

// GetTopicSubHandler 获取mqtt消息处理
func (td *Driver) GetTopicSubHandler() *map[string]mqttclient.SubTopic {
	return td.topicToSub
}
