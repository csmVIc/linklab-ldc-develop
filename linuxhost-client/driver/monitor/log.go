package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/linuxhost-client/driver/linuxhostnode"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) readlog(devport string) error {

	logchan, err := linuxhostnode.LDriver.GetDeviceReadLogChan(devport)
	if err != nil {
		log.Error("linuxhostnode.LDriver.GetDeviceReadLogChan {%v} error {%v}", devport, err)
		return err
	}

	for logmsg := range logchan {
		log.Debugf("read log {%v}", *logmsg)
		devicestatus, err := linuxhostnode.LDriver.GetDeviceStatus(devport)
		if err != nil {
			log.Error("linuxhostnode.LDriver.GetDeviceStatus devport {%v} error {%v}", devport, err)
			return err
		}
		if devicestatus.BusyStatus == linuxhostnode.Running {
			if err := topichandler.TDriver.PubDeviceLog(devicestatus.BurnInfo.GroupID, devicestatus.BurnInfo.TaskIndex, logmsg.Msg, logmsg.TimeStamp); err != nil {
				log.Error("topichandler.TDriver.PubDeviceLog {%v}", err)
			}
		}
	}

	return nil
}
