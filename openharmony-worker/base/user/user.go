package user

import "errors"

// Driver 接口调用
type Driver struct {
	info *UInfo
}

var (
	// UDriver 全局接口调用实例
	UDriver *Driver
)

// New 创建新的调用接口
func New(i *UInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil")
	}
	ud := &Driver{info: i}
	return ud, nil
}
