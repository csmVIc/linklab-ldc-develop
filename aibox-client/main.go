package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/aibox-client/driver/aiboxnode"
	"linklab/device-control-v2/aibox-client/driver/monitor"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/base-library/mqttclient"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		MqttInfo      mqttclient.MInfo   `json:"mqtt"`
		TopicInfo     topichandler.TInfo `json:"topic"`
		AIBoxNodeInfo aiboxnode.AInfo    `json:"aiboxnode"`
		MonitorInfo   monitor.MInfo      `json:"monitor"`
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
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error
	aiboxnode.ADriver, err = aiboxnode.New(&initInfo.AIBoxNodeInfo)
	if err != nil {
		log.Panicf("aiboxnode.New error {%s}", err)
	}
	// monitor初始化
	monitor.MDriver, err = monitor.New(&initInfo.MonitorInfo)
	if err != nil {
		log.Panicf("monitor.New error {%s}", err)
	}
	// mqtt消息订阅初始化
	mqttclient.MDriver, err = mqttclient.New(&initInfo.MqttInfo)
	if err != nil {
		log.Panicf("mqttclient.New error {%v}", err)
	}
	// mqtt消息句柄初始化
	topichandler.TDriver, err = topichandler.New(&initInfo.TopicInfo, monitor.MDriver.GetErrChan(), monitor.MDriver.GetBurnChan(), api.ADriver.GetTokenChan(), nil)
	if err != nil {
		log.Panicf("topichandler.New error {%v}", err)
	}
	go mqttInit()
	// monitor 运行
	if err := monitor.MDriver.Run(); err != nil {
		log.Panicf("monitor.MDriver.Run error {%v}", err)
	}
}

func mqttInit() {
	if err := mqttclient.MDriver.Monitor(topichandler.TDriver.GetTopicSubHandler()); err != nil {
		log.Panicf("mqttclient.MDriver.Monitor error {%v}", err)
	}
}
