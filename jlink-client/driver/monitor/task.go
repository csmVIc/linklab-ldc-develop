package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) taskprocess() error {

	defer func() {
		md.endrun = true
		log.Error("task func set endrun")
	}()

	for taskinfo := range md.taskchan {
		err := topichandler.TDriver.PubEndRun(taskinfo.BurnInfo, taskinfo.TimeOut, taskinfo.EndTime)
		if err != nil {
			log.Errorf("api.ADriver.RunEnd error {%v}", err)
			return err
		}
	}

	return nil
}

func (md *Driver) taskstartup() {
	for index := 0; index < runtime.NumCPU()*md.info.Task.ThreadMultiple; index++ {

		log.Debugf("task process {%v} start up", index)
		go md.taskprocess()
	}
}
