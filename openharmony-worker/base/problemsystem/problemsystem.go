package problemsystem

import "errors"

// Driver 判题接口
type Driver struct {
	info *PInfo
}

var (
	// PDriver 全局判题调用实例
	PDriver *Driver
)

// New 创建新的接口调用
func New(i *PInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}
	pd := &Driver{info: i}
	return pd, nil
}
