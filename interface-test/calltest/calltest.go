package calltest

import (
	"errors"
	"linklab/device-control-v2/interface-test/calltest/api"

	log "github.com/sirupsen/logrus"
)

// Driver 调用测试
type Driver struct {
	info    *CInfo
	endchan chan *TestThreadEnd
}

var (
	// CDriver 全局调用测试实例
	CDriver *Driver
)

// New 创建新的接口调用
func New(i *CInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}

	var err error
	api.ADriver, err = api.New(&i.API)
	if err != nil {
		log.Errorf("api.New error {%v}", err)
		return nil, err
	}

	cd := &Driver{info: i, endchan: make(chan *TestThreadEnd, i.EndChan.Size)}
	return cd, nil
}
