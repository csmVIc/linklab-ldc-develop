package api

import "errors"

// Driver 接口调用
type Driver struct {
	info *AInfo
}

var (
	// ADriver 全局接口调用实例
	ADriver *Driver
)

// New 创建新的接口调用
func New(i *AInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}
	ad := &Driver{info: i}
	return ad, nil
}
