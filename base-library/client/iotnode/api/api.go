package api

import "errors"

// Driver 接口调用
type Driver struct {
	info      *AInfo
	token     string
	tokenchan chan string
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
	ad := &Driver{info: i, tokenchan: make(chan string, i.Token.ChanSize)}
	return ad, nil
}
