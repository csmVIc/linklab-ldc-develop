package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"linklab/device-control-v2/project-deploy/driver"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	initInfo struct {
		ProjectList   []string `json:"project list"`
		RootDirectory string   `json:"root directory"`
		KubeConfig    string   `json:"kube config"`
		DockerConfig  string   `json:"docker config"`
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
	for index, projectname := range initInfo.ProjectList {
		log.Infof("[%d] project {%s} begin", index, projectname)
		if err := driver.ProjectDeploy(initInfo.RootDirectory, initInfo.KubeConfig, initInfo.DockerConfig, projectname); err != nil {
			log.Panicf("driver.ProjectDeploy error {%v}", err)
		}
		log.Infof("[%d] project {%s} begin", index, projectname)
	}
}
