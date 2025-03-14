package monitor

import (
	"linklab/device-control-v2/jlink-client/driver/iotnode"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) cmdwriteprocess() {
	for devicecmd := range md.cmdchan {
		err := iotnode.IDriver.DeviceCmdWrite(devicecmd)
		if err != nil {
			log.Error("device cmd {%v} write error {%v}", *devicecmd, err)
		}
	}
}

func (md *Driver) cmdwritestartup() {
	for index := 0; index < runtime.NumCPU()*md.info.CmdWrite.ThreadMultiple; index++ {

		log.Debugf("cmd write {%v} start up", index)
		go md.cmdwriteprocess()
	}
}
