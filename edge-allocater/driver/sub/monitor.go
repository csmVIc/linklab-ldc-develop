package sub

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Monitor 运行监控
func (sd *Driver) Monitor() error {
	for count := 0; count < sd.info.MaxReconn; count++ {
		err := sd.subscribe()
		if err != nil {
			log.Errorf("sd.subscriber() error {%v}", err)
		}
		time.Sleep(time.Duration(sd.info.ReconnInterval) * time.Second)
	}
	err := fmt.Errorf("sd.subscriber() error count >= limit {%v}", sd.info.MaxReconn)
	return err
}
