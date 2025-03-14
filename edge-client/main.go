package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/edge-client/driver/edgenode"
	"linklab/device-control-v2/edge-client/driver/monitor"
	"linklab/device-control-v2/edge-client/eapi"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		ServerInfo   eapi.EInfo         `json:"server"`
		MqttInfo     mqttclient.MInfo   `json:"mqtt"`
		TopicInfo    topichandler.TInfo `json:"topic"`
		MonitorInfo  monitor.MInfo      `json:"monitor"`
		EdgeNodeInfo edgenode.EInfo     `json:"edgenode"`
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
	// 随机数
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error
	// edgenode初始化
	edgenode.EDriver, err = edgenode.New(&initInfo.EdgeNodeInfo)
	if err != nil {
		log.Panicf("edgenode.New error {%s}", err)
	}
	// monitor初始化
	monitor.MDriver, err = monitor.New(&initInfo.MonitorInfo)
	if err != nil {
		log.Panicf("monitor.New error {%s}", err)
	}
	// mqtt消息订阅初始化
	log.Debugln(initInfo.MqttInfo)
	mqttclient.MDriver, err = mqttclient.New(&initInfo.MqttInfo)
	if err != nil {
		log.Panicf("mqttclient.New error {%v}", err)
	}
	// mqtt消息句柄初始化
	topichandler.TDriver, err = topichandler.New(&initInfo.TopicInfo)
	topichandler.TDriver.SetErrChan(monitor.MDriver.GetErrChan())
	topichandler.TDriver.SetTokenChan(api.ADriver.GetTokenChan())
	// topichandler.TDriver.SetPodApplyChan(monitor.MDriver.GetPodApplyChan())
	if err != nil {
		log.Panicf("topichandler.New error {%v}", err)
	}
	// mqtt 运行
	go mqttInit()
	// monitor token 初始化
	if err := monitor.MDriver.TokenInit(); err != nil {
		log.Panicf("monitor.MDriver.TokenInit error {%v}", err)
	}
	// monitor setup 初始化
	if err := monitor.MDriver.SetupInit(); err != nil {
		log.Panicf("monitor.MDriver.SetupInit error {%v}", err)
	}
	// monitor 运行
	go monitorInit()
	// router 运行
	if err := routerInit(); err != nil {
		log.Panicf("routerInit error {%v}", err)
	}
}

func mqttInit() {
	if err := mqttclient.MDriver.Monitor(topichandler.TDriver.GetTopicSubHandler()); err != nil {
		log.Panicf("mqttclient.MDriver.Monitor error {%v}", err)
	}
}

func monitorInit() {
	if err := monitor.MDriver.Run(); err != nil {
		log.Panicf("monitor.MDriver.Run error {%v}", err)
	}
}

func routerInit() error {
	router := gin.Default()
	router.Use(auth.CORS())

	lf, err := os.Create("./log/gin.log")
	if err != nil {
		log.Errorf("os.Create file error {%v}", err)
		return err
	}
	gin.DefaultWriter = io.MultiWriter(lf, os.Stdout)

	if err := eapi.RouterInit(router, &initInfo.ServerInfo); err != nil {
		log.Error("eapi.RouterInit error {%v}", err)
		return err
	}

	if err := router.Run(initInfo.ServerInfo.Address.Host + ":" + initInfo.ServerInfo.Address.Port); err != nil {
		log.Errorf("router.Run error {%v}", err)
		return err
	}

	return nil
}
