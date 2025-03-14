package judge

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

// Driver 负责任务分配的决策任务
type Driver struct {
	info *JInfo
}

var (
	// JDriver 全局实例
	JDriver *Driver
)

// New 创建实例
func New(i *JInfo) (*Driver, error) {
	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}
	md := &Driver{info: i}
	return md, nil
}
