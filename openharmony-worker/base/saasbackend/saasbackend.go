package saasbackend

import (
	"errors"
)

// Driver Saas接口
type Driver struct {
	info *SInfo
}

var (
	// SDriver 全局Saas调用实例
	SDriver *Driver
)

// New 创建新的接口调用
func New(i *SInfo) (*Driver, error) {

	if i == nil {
		return nil, errors.New("init info i nil error")
	}
	pd := &Driver{info: i}
	return pd, nil
}
