package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/virtual-device-client/driver/virtualnode"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) readlog(devport string) error {

	logchan, err := virtualnode.VDriver.ReadDeviceLog(devport, md.taskchan)
	if err != nil {
		log.Errorf("virtualnode.IDriver.ReadDeviceLog error {%v}", err)
		return err
	}

	for logmsg := range logchan {
		log.Debugf("read log {%v}", *logmsg)
		devicestatus, err := virtualnode.VDriver.GetDeviceStatus(devport)
		if err != nil {
			log.Error("virtualnode.IDriver.GetDeviceStatus devport {%v} error {%v}", devport, err)
			return err
		}
		if devicestatus.BusyStatus == virtualnode.Running {
			if err := topichandler.TDriver.PubDeviceLog(devicestatus.BurnInfo.GroupID, devicestatus.BurnInfo.TaskIndex, logmsg.Msg, logmsg.TimeStamp); err != nil {
				log.Error("topichandler.TDriver.PubDeviceLog {%v}", err)
			}
		}
	}

	return nil
}
