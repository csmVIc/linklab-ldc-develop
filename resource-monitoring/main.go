package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"linklab/device-control-v2/base-library/influx"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/resource-monitoring/monitor"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		InfluxInfo  influx.IInfo  `json:"influx"`
		MonitorInfo monitor.MInfo `json:"monitor"`
		LogInfo     logger.LInfo  `json:"log"`
	}
)

func init() {
	// log 初始化
	log.SetFormatter(&log.JSONFormatter{})
	lf, err := os.Create("./log/logrus.log")
	if err != nil {
		log.Panicf("os.Create log file error {%v}", err)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, lf))
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	// config 加载
	buf, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Panic(err)
	}
	// json 解析
	if err = json.Unmarshal(buf, &initInfo); err != nil {
		log.Panic(err)
	}
	// 随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	// influx初始化
	var err error
	influx.Idriver, err = influx.New(&initInfo.InfluxInfo)
	if err != nil {
		log.Panicf("influx.New error {%v}", err)
	}
	influx.Idriver.MonitorSetup()

	// log初始化
	logger.Ldriver, err = logger.New(&initInfo.LogInfo)
	if err != nil {
		log.Panicf("logger.New error {%s}", err)
	}

	// monitor init
	monitor.Mdriver, err = monitor.New(&initInfo.MonitorInfo)
	if err != nil {
		log.Panicf("monitor.New error {%v}", err)
	}
	err = monitor.Mdriver.Monitor()
	if err != nil {
		log.Panicf("monitor.Md.Monitor error {%v}", err)
	}
}
