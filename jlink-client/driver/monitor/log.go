package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/jlink-client/driver/iotnode"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) readlog(devport string) error {

	logchan, err := iotnode.IDriver.GetDeviceReadLogChan(devport)
	if err != nil {
		log.Errorf("iotnode.IDriver.GetDeviceReadLogChan error {%v}", err)
		return err
	}

	for logmsg := range logchan {
		log.Debugf("read log {%v}", *logmsg)
		devicestatus, err := iotnode.IDriver.GetDeviceStatus(devport)
		if err != nil {
			log.Error("iotnode.IDriver.GetDeviceStatus devport {%v} error {%v}", devport, err)
			return err
		}
		if devicestatus.BusyStatus == iotnode.Running {
			if err := topichandler.TDriver.PubDeviceLog(devicestatus.BurnInfo.GroupID, devicestatus.BurnInfo.TaskIndex, logmsg.Msg, logmsg.TimeStamp); err != nil {
				log.Error("topichandler.TDriver.PubDeviceLog {%v}", err)
			}
		}
	}

	return nil
}
