package aiboxnode

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// Driver 负责设备操作
type Driver struct {
	info    *AInfo
	devices Devices
}

// Driver 设备操作
var (
	ADriver *Driver
)

func (ad *Driver) init() error {
	if len(ad.info.Boards) < 1 {
		return errors.New("len(ad.info.Boards) < 1 error")
	}

	for key, value := range ad.info.Boards {
		if value.VirtualNum < 1 {
			return fmt.Errorf("board {%v} virtual number {%v} < 1", key, value.VirtualNum)
		}

		for index := 0; index < value.VirtualNum; index++ {
			ad.devices.devicesMap.Store(fmt.Sprintf("/dev/%s-%d", key, index), &DeviceStatus{
				BurnInfo:   nil,
				BeginTime:  time.Now(),
				BusyStatus: IdleState,
				LogChan:    make(chan *LogMsg, ad.info.DeviceLog.LogChanSize),
				Lock:       sync.Mutex{},
			})
		}
	}

	return nil
}

// New 创建设备操作
func New(i *AInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}

	ad := &Driver{info: i, devices: Devices{
		devicesMap:  sync.Map{},
		devicesLock: sync.RWMutex{},
	}}

	if err := ad.init(); err != nil {
		log.Errorf("ad.init() error {%v}", err)
		return nil, err
	}

	return ad, nil
}
