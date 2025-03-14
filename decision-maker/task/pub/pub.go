package pub

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Driver 负责消息发布
type Driver struct {
	info *PInfo
}

var (
	// PDriver 消息发布
	PDriver *Driver
)

// New 创建新实例
func New(i *PInfo) (*Driver, error) {
	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}
	pd := &Driver{info: i}
	return pd, nil
}
