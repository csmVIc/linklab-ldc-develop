package subscriber

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Driver 负责消息订阅
type Driver struct {
	info       *SInfo
	exitsignal bool
}

var (
	// MDriver 消息订阅全局实例
	MDriver *Driver
)

// New 创建执行实例
func New(i *SInfo) (*Driver, error) {
	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}
	sd := &Driver{info: i, exitsignal: false}
	return sd, nil
}
