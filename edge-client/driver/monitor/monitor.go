package monitor

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/api"
)

// Driver 监控
type Driver struct {
	info    *MInfo
	endrun  bool
	errchan chan error
	// podapplychan chan *msg.PodApply
}

// MDriver 监控全局实例
var (
	MDriver *Driver
)

// New 创建监控实例
func New(i *MInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}

	var err error
	api.ADriver, err = api.New(&i.API)
	if err != nil {
		return nil, fmt.Errorf("api.New error {%v}", err)
	}

	md := &Driver{info: i, endrun: false, errchan: make(chan error, 1)}
	return md, nil
}
