package virtualnode

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Driver 负责设备操作
type Driver struct {
	info    *VInfo
	devices Devices
}

// Driver 设备操作
var (
	VDriver *Driver
)

func (vd *Driver) init() error {
	if len(vd.info.Boards) < 1 {
		return errors.New("len(vd.info.Boards) < 1 error")
	}

	for key, value := range vd.info.Boards {
		if value.VirtualNum < 1 {
			return fmt.Errorf("board {%v} virtual number {%v} < 1", key, value.VirtualNum)
		}

		for index := 0; index < value.VirtualNum; index++ {
			vd.devices.devicesMap.Store(fmt.Sprintf("/dev/%s-%d", key, index), &DeviceStatus{
				BurnInfo:   nil,
				BeginTime:  time.Now(),
				BusyStatus: IdleState,
				Lock:       sync.Mutex{},
			})
		}
	}

	return nil
}

// New 创建设备操作
func New(i *VInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}

	vd := &Driver{info: i, devices: Devices{
		devicesMap:  sync.Map{},
		devicesLock: sync.RWMutex{},
	}}

	if err := vd.init(); err != nil {
		log.Errorf("vd.init() error {%v}", err)
		return nil, err
	}

	return vd, nil
}
