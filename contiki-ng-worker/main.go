package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/influx"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/contiki-ng-worker/driver/compile"
	"linklab/device-control-v2/contiki-ng-worker/driver/subscriber"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		MongoInfo      database.DInfo   `json:"mongo"`
		InfluxInfo     influx.IInfo     `json:"influx"`
		NatsInfo       messenger.MInfo  `json:"nats"`
		SubscriberInfo subscriber.SInfo `json:"subscriber"`
		LogInfo        logger.LInfo     `json:"log"`
		CompileInfo    compile.CInfo    `json:"compile"`
	}
)

func init() {
	// log 初始化
	log.SetFormatter(&log.JSONFormatter{})
	lf, err := os.Create("./log/logrus.log")
	if err != nil {
		log.Panicf("init os.Create log file error {%v}", err)
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
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error
	// influx初始化
	influx.Idriver, err = influx.New(&initInfo.InfluxInfo)
	if err != nil {
		log.Panicf("influx.New error {%v}", err)
	}
	influx.Idriver.MonitorSetup()
	// nats初始化
	messenger.Mdriver, err = messenger.New(&initInfo.NatsInfo)
	if err != nil {
		log.Panicf("messenger.New error {%v}", err)
	}
	// mongo初始化
	database.Mdriver, err = database.New(&initInfo.MongoInfo)
	if err != nil {
		log.Panicf("database.New error {%s}", err)
	}
	// log初始化
	logger.Ldriver, err = logger.New(&initInfo.LogInfo)
	if err != nil {
		log.Panicf("logger.New error {%s}", err)
	}
	// compile初始化
	compile.Cd, err = compile.New(&initInfo.CompileInfo)
	if err != nil {
		log.Panicf("compile.New error {%s}", err)
	}
	go compile.Cd.Monitor()
	// subscriber初始化
	subscriber.MDriver, err = subscriber.New(&initInfo.SubscriberInfo)
	if err != nil {
		log.Panicf("subscriber.New error {%s}", err)
	}
	if err = subscriber.MDriver.Monitor(); err != nil {
		log.Panicf("subscriber.MDriver.Monitor error {%v}", err)
	}
}
