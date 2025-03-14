package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/linuxhost-client/driver/linuxhostnode"
	"linklab/device-control-v2/linuxhost-client/driver/monitor"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		MqttInfo          mqttclient.MInfo    `json:"mqtt"`
		TopicInfo         topichandler.TInfo  `json:"topic"`
		LinuxHostNodeInfo linuxhostnode.LInfo `json:"linuxhostnode"`
		MonitorInfo       monitor.MInfo       `json:"monitor"`
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
	linuxhostnode.LDriver, err = linuxhostnode.New(&initInfo.LinuxHostNodeInfo)
	if err != nil {
		log.Panicf("linuxhostnode.New error {%s}", err)
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
	topichandler.TDriver, err = topichandler.New(&initInfo.TopicInfo)
	topichandler.TDriver.SetErrChan(monitor.MDriver.GetErrChan())
	topichandler.TDriver.SetBurnChan(monitor.MDriver.GetBurnChan())
	topichandler.TDriver.SetTokenChan(api.ADriver.GetTokenChan())
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
