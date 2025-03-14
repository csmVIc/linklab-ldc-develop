package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/edge-allocater/driver/sub"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		RedisInfo cache.CInfo     `json:"redis"`
		NatsInfo  messenger.MInfo `json:"nats"`
		MongoInfo database.DInfo  `json:"mongo"`
		SubInfo   sub.SInfo       `json:"sub"`
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
	var err error
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

	// wait 等待
	sub.SDriver, err = sub.New(&initInfo.SubInfo)
	if err != nil {
		log.Panicf("sub.New error {%v}", err)
	}
	if err = sub.SDriver.Monitor(); err != nil {
		log.Panicf("sub.SDriver.Monitor error {%v}", err)
	}
}
