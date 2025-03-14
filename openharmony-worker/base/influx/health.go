package influx

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
)

func (id *Driver) health() error {
	healthcheck, err := (*id.iclient).Health(context.TODO())
	if err != nil {
		log.Errorf("(*id.iclient).Health error {%v}", err)
		return err
	}
	if healthcheck.Status != "pass" {
		err := fmt.Errorf("(*id.iclient).Health.Status error {%v}", healthcheck.Status)
		log.Error(err)
		atomic.StoreInt32(&id.isHealth, 0)
		return err
	}
	atomic.StoreInt32(&id.isHealth, 1)
	return nil
}

func (id *Driver) healthCheck() {
	for {
		if err := id.health(); err != nil {
			log.Errorf("id.health() error {%v}", err)
		}
		time.Sleep(time.Second * time.Duration(id.info.Client.HealthCheckInterval))
	}
}

func (id *Driver) getHealthStatus() bool {
	value := atomic.LoadInt32(&id.isHealth)
	return value == 1
}
