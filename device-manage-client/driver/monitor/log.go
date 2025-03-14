package monitor

import (
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/device-manage-client/driver/iotnode"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) readlog(devport string) error {

	logchan, err := iotnode.IDriver.ReadDeviceLog(devport, md.taskchan)
	if err != nil {
		log.Errorf("iotnode.IDriver.ReadDeviceLog error {%v}", err)
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
			if logmsg == nil {
				// 可以在这里处理空指针情况，例如打印错误信息并跳过该条日志
				log.Error("Received nil log message")
				continue
			}
			
			// 在这里处理非空的 logmsg
			if devicestatus == nil {
				log.Error("devicestatus is nil")
				continue
			}
			
			if devicestatus.BurnInfo == nil {
				log.Error("BurnInfo is nil")
				continue
			}
			
			// 在这里处理非空的 devicestatus 和 devicestatus.BurnInfo
			log.Debugf("GroupID: {%v}", devicestatus.BurnInfo.GroupID)
			log.Debugf("TaskIndex: {%v}", devicestatus.BurnInfo.TaskIndex)
			
			if logmsg.Msg == "" {
				log.Error("logmsg.Msg is empty")
			} else {
				log.Debugf("Msg: {%v}", logmsg.Msg)
			}
			
			log.Debugf("TimeStamp: {%v}", logmsg.TimeStamp)
			if err := topichandler.TDriver.PubDeviceLog(devicestatus.BurnInfo.GroupID, devicestatus.BurnInfo.TaskIndex, logmsg.Msg, logmsg.TimeStamp); err != nil {
				log.Error("topichandler.TDriver.PubDeviceLog {%v}", err)
			}
		}
	}

	return nil
}
