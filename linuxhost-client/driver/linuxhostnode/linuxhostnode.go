package linuxhostnode

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Driver 负责设备操作
type Driver struct {
	info    *LInfo
	devices Devices
}

// Driver 设备操作
var (
	LDriver *Driver
)

func (ld *Driver) init() error {
	if len(ld.info.Boards) < 1 {
		return errors.New("len(ld.info.Boards) < 1 error")
	}

	for key, value := range ld.info.Boards {
		if value.VirtualNum < 1 {
			return fmt.Errorf("board {%v} virtual number {%v} < 1", key, value.VirtualNum)
		}

		for index := 0; index < value.VirtualNum; index++ {
			ld.devices.devicesMap.Store(fmt.Sprintf("/dev/%s-%d", key, index), &DeviceStatus{
				BurnInfo:   nil,
				BeginTime:  time.Now(),
				BusyStatus: IdleState,
				LogChan:    make(chan *LogMsg, ld.info.DeviceLog.LogChanSize),
				Lock:       sync.Mutex{},
			})
		}
	}

	return nil
}

// New 创建设备操作
func New(i *LInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}

	ld := &Driver{info: i, devices: Devices{
		devicesMap:  sync.Map{},
		devicesLock: sync.RWMutex{},
	}}

	if err := ld.init(); err != nil {
		log.Errorf("ld.init() error {%v}", err)
		return nil, err
	}

	return ld, nil
}
