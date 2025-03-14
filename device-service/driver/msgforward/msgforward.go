package msgforward

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Driver 负责客户端消息转发
type Driver struct {
	info       *MInfo
	exitsignal bool
}

var (
	// MDriver 消息转发全局实例
	MDriver *Driver
)

// New 创建执行实例
func New(i *MInfo) (*Driver, error) {
	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}
	md := &Driver{info: i, exitsignal: false}
	return md, nil
}
