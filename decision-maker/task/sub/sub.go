package sub

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Driver 负责订阅烧写任务
type Driver struct {
	info        *SInfo
	callbackErr error
}

var (
	// SDriver 订阅烧写任务全局实例
	SDriver *Driver
)

// New 创建执行实例
func New(i *SInfo) (*Driver, error) {
	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}
	md := &Driver{info: i, callbackErr: nil}
	return md, nil
}
