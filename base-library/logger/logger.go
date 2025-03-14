package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

// Driver 日志记录器
type Driver struct {
	info      *LInfo
	loggerMap map[string]*log.Logger
}

var (
	// Ldriver 全局实例
	Ldriver *Driver
)

// New 创建实例
func New(i *LInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info nil error")
	}
	log.Debugf("logger init info {%v}", i.Logs)
	ld := &Driver{info: i, loggerMap: make(map[string]*log.Logger)}
	for name, loginfo := range i.Logs {
		if len(loginfo.Directory) < 1 {
			ld.loggerMap[name] = nil
		} else {
			lf, err := os.Create(loginfo.Directory)
			if err != nil {
				return nil, err
			}
			ld.loggerMap[name] = log.New()
			ld.loggerMap[name].SetFormatter(&log.JSONFormatter{})
			ld.loggerMap[name].SetOutput(io.MultiWriter(os.Stdout, lf))
			ld.loggerMap[name].SetLevel(log.DebugLevel)
			ld.loggerMap[name].SetReportCaller(true)
			log.Debugf("logger create local logger {%v}", name)
		}
	}
	return ld, nil
}

// GetPointName 获取日志记录器的PointName
func (ld *Driver) getPointName(name string) (string, error) {
	if _, isOk := ld.info.Logs[name]; isOk == false {
		err := fmt.Errorf("info doesn't contain {%v} error", name)
		return "", err
	}
	return ld.info.Logs[name].Pointname, nil
}

// GetLogger 获取日志记录器的本地Logger
func (ld *Driver) getLogger(name string) (*log.Logger, error) {
	if _, isOk := ld.loggerMap[name]; isOk == false {
		err := fmt.Errorf("logger map doesn't contain {%v} error", name)
		return nil, err
	}
	return ld.loggerMap[name], nil
}
