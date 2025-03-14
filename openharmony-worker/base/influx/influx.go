package influx

import (
	"errors"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// Driver influx客户端
type Driver struct {
	info     *IInfo
	iclient  *influxdb2.Client
	PushChan chan *Point
	isHealth int32
}

var (
	// Idriver influx客户端全局实例
	Idriver *Driver
)

func (id *Driver) init() {
	c := influxdb2.NewClientWithOptions(id.info.Client.URL, fmt.Sprintf("%s:%s", id.info.Client.UserName, id.info.Client.PassWord), influxdb2.DefaultOptions().SetBatchSize(id.info.Client.BatchSize).SetFlushInterval(id.info.Client.FlushInterval).SetUseGZip(id.info.Client.UseGZip))
	id.iclient = &c
}

// New 创建执行实例
func New(i *IInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info nil error")
	}
	id := &Driver{info: i, iclient: nil, PushChan: make(chan *Point, i.Chans.Size), isHealth: 0}
	id.init()
	return id, nil
}
