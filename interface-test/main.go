package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/interface-test/calltest"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		CallTest calltest.CInfo `json:"calltest"`
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
	buf, err := ioutil.ReadFile("./config/k8s-config.json")
	if err != nil {
		log.Panic(err)
	}
	// json 解析
	if err = json.Unmarshal(buf, &initInfo); err != nil {
		log.Panic(err)
	}

	// 随机数
	rand.Seed(time.Now().Unix())
}

func main() {
	var err error
	// calltest初始化
	calltest.CDriver, err = calltest.New(&initInfo.CallTest)
	if err != nil {
		log.Panicf("calltest.New error {%v}", err)
	}
	// calltest运行
	// for requesttimes := initInfo.CallTest.Test.BeginTimes; requesttimes <= initInfo.CallTest.Test.EndTimes; requesttimes += initInfo.CallTest.Test.Step {
	// 	err = calltest.CDriver.Run(requesttimes)
	// 	if err != nil {
	// 		log.Panicf("calltest.CDriver.Run error {%v}", err)
	// 	}
	// 	// 等待3分钟
	// 	time.Sleep(time.Minute * 10)
	// }

	for {
		err = calltest.CDriver.RunGroup()
		if err != nil {
			log.Panicf("calltest.CDriver.RunGroup error {%v}", err)
		}
		time.Sleep(time.Minute * 10)

		for i := 0; i < len(initInfo.CallTest.Test.Groups); i++ {

			initInfo.CallTest.Test.Groups[i].TotalTimes += initInfo.CallTest.Test.Groups[i].Step

			if initInfo.CallTest.Test.Groups[i].TotalTimes > initInfo.CallTest.Test.Groups[i].EndTimes {
				return
			}

		}
	}

}
