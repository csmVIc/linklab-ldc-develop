package influx

import (
	"fmt"
	"runtime"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

// MonitorSetup 启动point的监控协程
func (id *Driver) MonitorSetup() {
	go id.healthCheck()
	for index := 0; index < runtime.NumCPU()*id.info.Chans.ThreadMultiple; index++ {
		go id.taskMonitor()
	}
	log.Infof("MonitorSetup create monitor")
}

func (id *Driver) taskMonitor() {
	writeAPI := (*id.iclient).WriteAPI("", fmt.Sprintf("%s/autogen", id.info.Client.DataBase))

	// 错误处理
	errorsChan := writeAPI.Errors()
	go func() {
		for err := range errorsChan {
			log.Errorf("writeAPI.Errors {%v}", err)
		}
	}()

	for {
		select {
		case elem := <-id.PushChan:
			go id.writePoint(&writeAPI, elem)
		}
	}
}

func (id *Driver) writePoint(writeAPI *influxdb2api.WriteAPI, elem *Point) {
	point := influxdb2.NewPoint(elem.Measurement, elem.Tags, elem.Fields, time.Now())
	(*writeAPI).WritePoint(point)
}
