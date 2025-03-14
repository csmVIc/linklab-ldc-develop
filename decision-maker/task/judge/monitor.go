package judge

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Monitor 监听任务
func (jd *Driver) Monitor() error {
	go jd.monitor()
	go jd.groupMonitor()
	return nil
}

func (jd *Driver) monitor() {
	for count := 0; count < jd.info.MaxReconn; count++ {
		err := jd.listen()
		if err != nil {
			log.Errorf("jd.listen() error {%v}", err)
		}
		time.Sleep(time.Duration(jd.info.ReconnInterval) * time.Second)
	}
	log.Panicf("jd.listen() error count >= limit {%v}", jd.info.MaxReconn)
}

func (jd *Driver) groupMonitor() {
	for count := 0; count < jd.info.MaxReconn; count++ {
		err := jd.groupListen()
		if err != nil {
			log.Errorf("jd.groupListen() error {%v}", err)
		}
		time.Sleep(time.Duration(jd.info.ReconnInterval) * time.Second)
	}
	log.Panicf("jd.groupListen() error count >= limit {%v}", jd.info.MaxReconn)
}
