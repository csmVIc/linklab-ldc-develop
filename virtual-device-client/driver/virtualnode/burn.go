package virtualnode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// Burn 烧写
func (vd *Driver) Burn(burninfo *msg.ClientBurnMsg, burnfile string) (string, error) {

	defer func() {
		// 不管烧写成功或者失败,烧写结束时需要将文件删除
		if err := os.Remove(burnfile); err != nil {
			log.Errorf("burnfile {%v} remove error", burnfile)
		}
	}()

	// 设备操作加锁
	value, isOk := vd.devices.devicesMap.Load(burninfo.DeviceID)
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
	if err = vd.setburn(burninfo.DeviceID); err != nil {
		err = fmt.Errorf("devport {%v} set burn state error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return "", err
	}

	// 如果烧写成功则设置设备为运行状态,否则为空闲状态
	defer func() {
		if err != nil {
			if err = vd.setidle(burninfo.DeviceID); err != nil {
				err = fmt.Errorf("devport {%v} set idle state error {%v}", burninfo.DeviceID, err)
				log.Error(err)
			}
		} else {
			if err = vd.setrun(burninfo.DeviceID); err != nil {
				err = fmt.Errorf("devport {%v} set run state error {%v}", burninfo.DeviceID, err)
				log.Error(err)
			}
		}
	}()

	// 设置烧写信息
	devicestatus.BurnInfo = burninfo

	// 解析出board
	board := ""
	if board, err = vd.GetBoardFromDevPort(burninfo.DeviceID); err != nil {
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

	// 模拟烧写延迟
	time.Sleep(time.Duration(vd.info.Boards[board].BurnDelay) * time.Second)

	return "burn success", nil
}
