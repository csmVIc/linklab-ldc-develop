package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/influx"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/file-cache/api"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		ServerInfo api.AInfo      `json:"server"`
		InfluxInfo influx.IInfo   `json:"influx"`
		RedisInfo  cache.CInfo    `json:"redis"`
		MongoInfo  database.DInfo `json:"mongo"`
		LogInfo    logger.LInfo   `json:"log"`
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
	// redis初始化
	cache.Cdriver, err = cache.New(&initInfo.RedisInfo)
	if err != nil {
		log.Panicf("cache.New error {%v}", err)
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
	// 路由初始化
	if err = routerInit(); err != nil {
		log.Panicf("routerInit error {%v}", err)
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

	if err := api.RouterInit(router, &initInfo.ServerInfo.FileCache); err != nil {
		log.Error("api.RouterInit error {%v}", err)
		return err
	}

	if err := router.Run(initInfo.ServerInfo.Address.Host + ":" + initInfo.ServerInfo.Address.Port); err != nil {
		log.Errorf("router.Run error {%v}", err)
		return err
	}

	return nil
}
