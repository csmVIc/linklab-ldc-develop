package mqttclient

import (
	"sync/atomic"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

// 默认的消息接收回调函数
// 如果没有指定其它的消息处理函数,则调用该函数
func (md *Driver) msghandler(client mqtt.Client, msg mqtt.Message) {
	log.Infof("mqtt recv msg topic {%s} body {%s}", msg.Topic(), string(msg.Payload()))
}

// 连接丢失回调函数
func (md *Driver) connlosthandler(client mqtt.Client, err error) {
	atomic.AddInt32(&md.disconnCount, 1)
	log.Errorf("mqtt conn lost error {%v}, disconn count {%v}", err, atomic.LoadInt32(&md.disconnCount))
}

// 建立连接回调函数
func (md *Driver) onconnhandler(c mqtt.Client) {
	log.Info("mqtt conn set up begin")
	atomic.StoreInt32(&md.disconnCount, 0)
	if md.topicToSub != nil {
		if err := md.subinit(); err != nil {
			log.Errorf("mqtt topic sub error {%v}", err)
			atomic.AddInt32(&md.subInitErrSignal, 1)
		}
	} else {
		log.Infof("mqtt topic to sub nil")
	}
}
