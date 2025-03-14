package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/compile-gateway/api"
)

var (
	initInfo struct {
		MongoInfo  database.DInfo  `json:"mongo"`
		ServerInfo api.AInfo       `json:"server"`
		NatsInfo   messenger.MInfo `json:"nats"`
	}
)

func init() {
	// log 初始化
	log.SetFormatter(&log.JSONFormatter{})
	lf, err := os.Create("./log/logrus.log")
	if err != nil {
		log.Panicf("Init os.Create log file error {%v}", err)
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
	// nats初始化
	var err error
	messenger.Mdriver, err = messenger.New(&initInfo.NatsInfo)
	if err != nil {
		log.Panicf("messenger.New error {%v}", err)
	}
	// mongo 初始化
	database.Mdriver, err = database.New(&initInfo.MongoInfo)
	if err != nil {
		log.Panicf("database.New error {%s}", err)
	}
	// router init
	routerInit()
}

func routerInit() error {
	router := gin.Default()
	router.Use(auth.CORS())

	// gin log
	lf, err := os.Create("./log/gin.log")
	if err != nil {
		log.Errorf("routerInit os.Create log file error {%v}", err)
		return err
	}
	gin.DefaultWriter = io.MultiWriter(lf, os.Stdout)

	if err := api.RouterInit(router, &initInfo.ServerInfo.Handler); err != nil {
		log.Error("api.RouterInit return false")
		return err
	}

	if err := router.Run(initInfo.ServerInfo.Address.Host + ":" + initInfo.ServerInfo.Address.Port); err != nil {
		log.Errorf("router.Run error {%v}", err)
		return err
	}

	return nil
}
