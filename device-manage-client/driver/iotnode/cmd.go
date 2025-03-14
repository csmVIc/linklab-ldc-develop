package iotnode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// DeviceCmdWrite 设备命令写入
func (id *Driver) DeviceCmdWrite(devicecmd *msg.DeviceCmd) error {

	// 解析出设备状态
	value, isExist := id.devices.devicesMap.Load(devicecmd.DeviceID)
	if isExist == false {
		err := fmt.Errorf("devport {%v} not exist error", devicecmd.DeviceID)
		log.Error(err)
		return err
	}
	devicestatus := value.(*DeviceStatus)

	// 检查设备的运行状态
	if devicestatus.BusyStatus != Burning && devicestatus.BusyStatus != Running {
		err := fmt.Errorf("devport {%v} not burning or running error", devicecmd.DeviceID)
		log.Error(err)
		return err
	}

	// 检查设备所属任务和命令ID是否相同
	burninfo := devicestatus.BurnInfo
	if burninfo.GroupID != devicecmd.GroupID || burninfo.TaskIndex != devicecmd.TaskIndex {
		err := fmt.Errorf("devport {%v} burninfo {%v:%v} not equal devicecmd {%v:%v}", devicecmd.DeviceID, burninfo.GroupID, burninfo.TaskIndex, devicecmd.GroupID, devicecmd.TaskIndex)
		log.Error(err)
		return err
	}

	// 将待写入命令加入到通道
	(*devicestatus.CmdChan) <- devicecmd
	return nil
}
