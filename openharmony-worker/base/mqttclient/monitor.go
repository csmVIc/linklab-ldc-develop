package mqttclient

import (
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

// Monitor mqtt运行监控
func (md *Driver) Monitor(ts *map[string]SubTopic) error {

	time.Sleep(time.Second)
	// 连接测试
	if err := md.ping(); err != nil {
		log.Errorf("mqtt client ping error {%v}", err)
		return err
	}

	// 初始化订阅函数
	md.topicToSub = ts
	if err := md.subinit(); err != nil {
		log.Errorf("mqtt client topic sub init error {%v}", err)
		return err
	}

	// 运行监控
	for int(atomic.LoadInt32(&md.disconnCount)) < md.info.Monitor.MaxDisconnWait {
		time.Sleep(time.Second)
		// 检测mqtt重新连接订阅过程中的错误
		if atomic.LoadInt32(&md.subInitErrSignal) > 0 {
			err := fmt.Errorf("capture mqtt client sub init error {%v}", atomic.LoadInt32(&md.subInitErrSignal))
			return err
		}
	}

	// 错误连接计数,达到最大计数值,退出运行监控
	err := fmt.Errorf("mqtt client disconn count {%v} >= max disconn count {%v} error", atomic.LoadInt32(&md.disconnCount), md.info.Monitor.MaxDisconnWait)
	log.Error(err)
	return err
}
