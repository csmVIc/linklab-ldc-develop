package influx

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// PushPoint Point加入到处理队列
func (id *Driver) PushPoint(p *Point) error {
	if p == nil {
		err := errors.New("p is nil error")
		log.Error(err)
		return err
	}
	if id.getHealthStatus() == false {
		err := errors.New("id.getHealthStatus() == false")
		log.Error(err)
		return err
	}
	timeout := time.Duration(id.info.Chans.TimeOut) * time.Second
	select {
	case id.PushChan <- p:
	case <-time.After(timeout):
		err := fmt.Errorf("id.PushChan <- p timeout %vs", timeout)
		log.Error(err)
		return err
	}
	return nil
}
