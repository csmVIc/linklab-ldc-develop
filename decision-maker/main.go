package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/influx"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/decision-maker/task"
	"linklab/device-control-v2/decision-maker/task/judge"
	"linklab/device-control-v2/decision-maker/task/pub"
	"linklab/device-control-v2/decision-maker/task/sub"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		InfluxInfo influx.IInfo    `json:"influx"`
		RedisInfo  cache.CInfo     `json:"redis"`
		NatsInfo   messenger.MInfo `json:"nats"`
		MongoInfo  database.DInfo  `json:"mongo"`
		TaskInfo   task.TInfo      `json:"task"`
		LogInfo    logger.LInfo    `json:"log"`
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
	// 生成随机数种子
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
	// redis初始化
	cache.Cdriver, err = cache.New(&initInfo.RedisInfo)
	if err != nil {
		log.Panicf("cache.New error {%v}", err)
	}
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
	// sub初始化
	sub.SDriver, err = sub.New(&initInfo.TaskInfo.Sub)
	if err != nil {
		log.Panicf("sub.New error {%s}", err)
	}
	err = sub.SDriver.Monitor()
	if err != nil {
		log.Panicf("sub.SDriver.Monitor error {%s}", err)
	}
	// pub 初始化
	pub.PDriver, err = pub.New(&initInfo.TaskInfo.Pub)
	if err != nil {
		log.Panicf("pub.New error {%s}", err)
	}
	// judge 初始化
	judge.JDriver, err = judge.New(&initInfo.TaskInfo.Judge)
	if err != nil {
		log.Panicf("judge.New error {%s}", err)
	}
	err = judge.JDriver.Monitor()
	if err != nil {
		log.Panicf("judge.JDriver.Listen error {%s}", err)
	}
	// log初始化
	logger.Ldriver, err = logger.New(&initInfo.LogInfo)
	if err != nil {
		log.Panicf("logger.New error {%s}", err)
	}
	// wait 等待
	for {
		time.Sleep(time.Second)
	}
}
