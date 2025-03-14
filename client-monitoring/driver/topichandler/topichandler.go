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
	ts["$share/client-monitoring/$SYS/brokers/+/clients/+/connected"] = mqttclient.SubTopic{
		MsgHandler: td.clientconnect,
		Qos:        0,
	}
	ts["$share/client-monitoring/$SYS/brokers/+/clients/+/disconnected"] = mqttclient.SubTopic{
		MsgHandler: td.clientdisconnect,
		Qos:        0,
	}
	ts["$share/client-monitoring/clients/+/+/publish/client/heartbeat"] = mqttclient.SubTopic{
		MsgHandler: td.clientheartbeat,
		Qos:        2,
	}
	ts["$share/client-monitoring/clients/+/+/publish/device/update"] = mqttclient.SubTopic{
		MsgHandler: td.deviceupdate,
		Qos:        2,
	}
	ts["$share/client-monitoring/clients/+/+/publish/edgenode/update"] = mqttclient.SubTopic{
		MsgHandler: td.edgenodeupdate,
		Qos:        2,
	}
	ts["$share/client-monitoring/clients/+/+/publish/pod/update"] = mqttclient.SubTopic{
		MsgHandler: td.podupdate,
		Qos:        2,
	}
	ts["$share/client-monitoring/clients/+/+/publish/edgenode/resource"] = mqttclient.SubTopic{
		MsgHandler: td.edgenoderesource,
		Qos:        2,
	}
	ts["$share/client-monitoring/clients/+/+/publish/pod/resource"] = mqttclient.SubTopic{
		MsgHandler: td.podresource,
		Qos:        2,
	}
	ts["$share/client-monitoring/clients/+/+/publish/edgenode/setup"] = mqttclient.SubTopic{
		MsgHandler: td.edgenodesetup,
		Qos:        2,
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
