package logger

import (
	"fmt"
	"linklab/device-control-v2/base-library/influx"

	log "github.com/sirupsen/logrus"
)

// WriteLog 记录日志
func (ld *Driver) WriteLog(name string, tags map[string]string, fields map[string]interface{}) error {
	// 本地日志记录
	locallogger, err := ld.getLogger(name)
	if err != nil {
		err = fmt.Errorf("getLogger err {%v}", err)
		log.Error(err)
		return err
	}
	if locallogger != nil {
		localfields := make(map[string]interface{})
		for key, valye := range tags {
			localfields[key] = valye
		}
		for key, value := range fields {
			localfields[key] = value
		}
		locallogger.WithFields(localfields).Info("write log")
	}

	// 数据库日志记录
	pointname, err := ld.getPointName(name)
	if err != nil {
		err = fmt.Errorf("getPointName err {%v}", err)
		log.Error(err)
		return err
	}
	point := &influx.Point{
		Measurement: pointname,
		Tags:        tags,
		Fields:      fields,
	}
	err = influx.Idriver.PushPoint(point)
	if err != nil {
		err = fmt.Errorf("influx.Idriver.PushPoint err {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}
