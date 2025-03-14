package monitor

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Monitor 依次运行相应的资源监控协程
func (md *Driver) Monitor() error {
	if err := md.check(); err != nil {
		log.Errorf("check error {%v}", err)
		return err
	}
	md.exitwg.Add(1)
	go md.podmonitor()
	go md.nodemonitor()
	md.exitwg.Wait()
	return nil
}

// 监控前的检查
func (md *Driver) check() error {
	if md.clientset == nil && md.k8sinit() != nil {
		err := errors.New("clientset k8sinit error")
		log.Error(err)
		return err
	}
	return nil
}
