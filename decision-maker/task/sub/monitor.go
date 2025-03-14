package sub

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Monitor 启动订阅任务监控协程
func (sd *Driver) Monitor() error {
	go sd.monitor()
	return nil
}

func (sd *Driver) monitor() {
	for count := 0; count < sd.info.MaxReconn; count++ {
		err := sd.subscriber()
		if err != nil {
			log.Errorf("sd.subscriber() error {%v}", err)
		}
		time.Sleep(time.Duration(sd.info.ReconnInterval) * time.Second)
	}
	log.Panicf("sd.subscriber() error count >= limit {%v}", sd.info.MaxReconn)
}
