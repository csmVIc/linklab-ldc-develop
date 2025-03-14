package influxv1

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client"
	client "github.com/influxdata/influxdb1-client/v2"
)

// Driver influx客户端
type Driver struct {
	info    *IInfo
	iclient *client.Client
}

var (
	// IDriver influx客户端全局实例
	IDriver *Driver
)

func (id *Driver) ping() error {
	_, _, err := (*id.iclient).Ping(time.Second * time.Duration(id.info.Client.PingTimeOut))
	return err
}

func (id *Driver) init() error {
	iclient, err := client.NewHTTPClient(
		client.HTTPConfig{
			Addr:     id.info.Client.URL,
			Username: id.info.Client.UserName,
			Password: id.info.Client.PassWord,
			Timeout:  time.Second * time.Duration(id.info.Client.WriteTimeOut),
		},
	)
	if err != nil {
		return fmt.Errorf("client.NewHTTPClient error {%v}", err)
	}
	id.iclient = &iclient
	return nil
}

// New 创建执行实例
func New(i *IInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info nil error")
	}
	id := &Driver{info: i, iclient: nil}
	if err := id.init(); err != nil {
		return nil, fmt.Errorf("id.init error {%v}", err)
	}
	return id, nil
}
