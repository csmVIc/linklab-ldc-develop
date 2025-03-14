package monitor

import (
	"linklab/device-control-v2/aibox-client/driver/aiboxnode"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) readlog(devport string) error {

	logchan, err := aiboxnode.ADriver.GetDeviceReadLogChan(devport)
	if err != nil {
		log.Error("aiboxnode.ADriver.GetDeviceReadLogChan {%v} error {%v}", devport, err)
		return err
	}

	for logmsg := range logchan {
		log.Debugf("read log {%v}", *logmsg)
		devicestatus, err := aiboxnode.ADriver.GetDeviceStatus(devport)
		if err != nil {
			log.Error("aiboxnode.ADriver.GetDeviceStatus devport {%v} error {%v}", devport, err)
			return err
		}
		if devicestatus.BusyStatus == aiboxnode.Running {
			if err := topichandler.TDriver.PubDeviceLog(devicestatus.BurnInfo.GroupID, devicestatus.BurnInfo.TaskIndex, logmsg.Msg, logmsg.TimeStamp); err != nil {
				log.Error("topichandler.TDriver.PubDeviceLog {%v}", err)
			}
		}
	}

	return nil
}
