package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/influx"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/client-monitoring/driver/topichandler"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		MqttInfo   mqttclient.MInfo   `json:"mqtt"`
		InfluxInfo influx.IInfo       `json:"influx"`
		RedisInfo  cache.CInfo        `json:"redis"`
		NatsInfo   messenger.MInfo    `json:"nats"`
		MongoInfo  database.DInfo     `json:"mongo"`
		TopicInfo  topichandler.TInfo `json:"topic"`
		LogInfo    logger.LInfo       `json:"log"`
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
	// log初始化
	logger.Ldriver, err = logger.New(&initInfo.LogInfo)
	if err != nil {
		log.Panicf("logger.New error {%s}", err)
	}
	// mqtt消息处理初始化
	topichandler.TDriver, err = topichandler.New(&initInfo.TopicInfo)
	if err != nil {
		log.Panicf("topichandler.New error {%v}", err)
	}
	// mqtt消息订阅初始化
	mqttclient.MDriver, err = mqttclient.New(&initInfo.MqttInfo)
	if err != nil {
		log.Panicf("mqttclient.New error {%v}", err)
	}
	err = mqttclient.MDriver.Monitor(topichandler.TDriver.GetTopicSubHandler())
	if err != nil {
		log.Panicf("mqttclient.MDriver.Monitor error {%v}", err)
	}
}
