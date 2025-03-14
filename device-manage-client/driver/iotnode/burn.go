package iotnode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"

	log "github.com/sirupsen/logrus"
)

// Burn 烧写
func (id *Driver) Burn(burninfo *msg.ClientBurnMsg, burnfile string) (string, error) {

	defer func() {
		// 不管烧写成功或者失败,烧写结束时需要将文件删除
		if err := os.Remove(burnfile); err != nil {
			log.Errorf("burnfile {%v} remove error", burnfile)
		}
	}()

	// 设备操作加锁
	value, isOk := id.devices.devicesMap.Load(burninfo.DeviceID)
	if isOk == false {
		err := fmt.Errorf("devport {%v} not in devicesMap error", burninfo.DeviceID)
		log.Error(err)
		return "", nil
	}
	devicestatus := value.(*DeviceStatus)
	devicestatus.Lock.Lock()
	defer func() {
		devicestatus.Lock.Unlock()
	}()

	var err error = nil

	// 如果串口存在则设置正在烧写
	if err = id.setburn(burninfo.DeviceID); err != nil {
		err = fmt.Errorf("devport {%v} set burn state error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return "", err
	}

	// 如果烧写成功则设置设备为运行状态,否则为空闲状态
	defer func() {
		if err != nil {
			if err = id.setidle(burninfo.DeviceID); err != nil {
				err = fmt.Errorf("devport {%v} set idle state error {%v}", burninfo.DeviceID, err)
				log.Error(err)
			}
		} else {
			if err = id.setrun(burninfo.DeviceID); err != nil {
				err = fmt.Errorf("devport {%v} set run state error {%v}", burninfo.DeviceID, err)
				log.Error(err)
			}
		}
	}()

	// 设置烧写信息
	devicestatus.BurnInfo = burninfo

	// 解析出board
	board := ""
	if board, err = id.GetBoardFromDevPort(burninfo.DeviceID); err != nil {
		err = fmt.Errorf("get board from devport error {%v}", err)
		log.Error(err)
		return "", err
	}

	// 检查文件是否存在
	if _, err := os.Stat(burnfile); os.IsNotExist(err) {
		err = fmt.Errorf("burnfile {%v} not exist", burnfile)
		log.Error(err)
		return "", err
	}

	// 烧写重试
	var burnmsg string = ""
	for index := 0; index < id.info.Burn.MaxRetryTimes; index++ {
		// 烧写
		burnmsg, err = id.boardCmdMap[board].Burn(burninfo.DeviceID, burnfile, id.info.Commands[board].Burn, devicestatus.SerialPort)
		log.Debugf("device burn msg\n%v", burnmsg)
		if err != nil {
			err = fmt.Errorf("device {%v} burn error {%v}", burninfo.DeviceID, err)
			log.Error(err)
			continue
		}

		return burnmsg, nil
	}

	return burnmsg, err
}
