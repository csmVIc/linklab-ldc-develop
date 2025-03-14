package subscriber

import (
	"errors"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/alios-edu-worker/driver/compile"
	"time"

	log "github.com/sirupsen/logrus"
)

// Monitor 接收订阅消息
func (sd *Driver) Monitor() error {

	compiletypes := compile.Cd.GetSupportCompileType()
	if compiletypes == nil || len(compiletypes) < 1 {
		err := errors.New("support compile type is nil or len < 1 error")
		log.Error(err)
		return err
	}

	for _, compiletype := range compiletypes {
		go sd.compilehandler(compiletype, compile.Cd.CheckCompileTypeSupportSystem(compiletype))
	}

	for {
		if isOk := messenger.Mdriver.GetClosed(); isOk == true {
			return errors.New("messenger.Mdriver.GetClosed true error")
		}
		if sd.exitsignal == true {
			return errors.New("sub handler exit error")
		}
		time.Sleep(time.Second)
	}
}
