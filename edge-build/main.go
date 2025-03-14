package main

import (
	"io"
	"math/rand"
	"os"
	"time"
	"/Users/csmvic/Desktop/linklab-ldc-develop/edge-build/driver"
	log "github.com/sirupsen/logrus"
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
	// 随机数
	rand.Seed(time.Now().UnixNano())
}

func main() {
	var err error
	driver.DDriver, err = driver.New()
	if err != nil {
		log.Panicf("driver.New error {%v}", err)
	}

	driver.DDriver.Execute()
}
