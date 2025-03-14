package calltest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/interface-test/calltest/api"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func (cd *Driver) testthread(datadir string, tmpdir string, index int) {

	// 用户名
	username := fmt.Sprintf("%s-%v", cd.info.Login.UserName, index)

	// 结束时发送退出信息
	endinfo := &TestThreadEnd{
		UserName: username,
		Err:      nil,
	}
	defer func() {
		cd.endchan <- endinfo
	}()

	// 创建日志记录map
	logMap := make(map[string]interface{})
	logMap["username"] = username

	// 判断是否需要编译
	compileoutput := []byte{}
	notexist := false
	var err error = nil
	if cd.info.Compile.NeedCompile {

		// 创建该线程的操作目录
		workdir := filepath.Join(tmpdir, username)
		if err := os.Mkdir(workdir, 0755); err != nil {
			err := fmt.Errorf("os.Mkdir {%v} error {%v}", workdir, err)
			log.Error(err)
			endinfo.Err = err
			return
		}

		// 获取编译源文件
		compilezip, err := cd.generateSource(workdir, index)
		if err != nil {
			err := fmt.Errorf("cd.generateSource error {%v}", err)
			log.Error(err)
			endinfo.Err = err
			return
		}

		// 未查询到结果,重新提交
	reupload:

		// 提交编译
		logMap["begincompile"] = time.Now().UnixNano()
		compilefilehash := ""
		for index := 0; index < cd.info.Burn.RetryTimes; index++ {
			compilefilehash, err = api.ADriver.CompileUpload(cd.info.Compile.BoardType, cd.info.Compile.CompileType, compilezip)
			if err != nil {
				err := fmt.Errorf("api.ADriver.CompileUpload username {%v} error {%v}", username, err)
				log.Error(err)
			} else {
				break
			}
		}
		if err != nil {
			endinfo.Err = err
		}
		logMap["compilefilehash"] = compilefilehash

		// 下载编译结果
		for len(compileoutput) < 1 {
			compileoutput, notexist, err = api.ADriver.CompileDownload(cd.info.Compile.BoardType, cd.info.Compile.CompileType, compilefilehash)
			// if err != nil {
			// 	err := fmt.Errorf("api.ADriver.CompileDownload error {%v}", err)
			// 	log.Error(err)
			// 	endinfo.Err = err
			// 	return
			// }

			if len(compileoutput) < 1 && notexist {
				goto reupload
			}

			time.Sleep(time.Second)
			log.Debugln("download sleep ...")
		}

		logMap["endcompile"] = time.Now().UnixNano()
	}

	// 用户登录
	logMap["beginuserlogin"] = time.Now().UnixNano()
	token := ""
	for index := 0; index < cd.info.Burn.RetryTimes; index++ {
		token, err = api.ADriver.UserLogin(username, cd.info.Login.PassWord)
		if err != nil {
			err := fmt.Errorf("api.ADriver.UserLogin username {%v} error {%v}", username, err)
			log.Error(err)
		} else {
			break
		}
	}
	if err != nil {
		endinfo.Err = err
	}
	logMap["enduserlogin"] = time.Now().UnixNano()

	// 文件上传
	logMap["beginfileupload"] = time.Now().UnixNano()
	filehash := ""
	err = nil
	for index := 0; index < cd.info.Burn.RetryTimes; index++ {
		filehash, err = api.ADriver.FileUpload(cd.info.Burn.BoardName, cd.info.Burn.FileSize, cd.info.Burn.FileRandom, token, cd.info.Burn.FilePath, compileoutput)
		if err != nil {
			err := fmt.Errorf("api.ADriver.FileUpload username {%v} error {%v}", username, err)
			log.Error(err)
		} else {
			break
		}
	}
	if err != nil {
		endinfo.Err = err
	}
	logMap["endfileupload"] = time.Now().UnixNano()

	// 日志监听
	handle, err := api.ADriver.GetWebSocketHandle(token)
	if err != nil {
		err := fmt.Errorf("api.ADriver.GetWebSocketHandle username {%v} error {%v}", username, err)
		log.Error(err)
		endinfo.Err = err
		return
	}
	endchan := make(chan error)
	go cd.logmonitor(handle, &logMap, endchan)

	// 烧写任务
	parameter := &request.DeviceBurnTasks{
		Tasks: []request.DeviceBurnTask{
			{
				BoardName: cd.info.Burn.BoardName,
				DeviceID:  "",
				RunTime:   cd.info.Burn.RunTime,
				FileHash:  filehash,
				ClientID:  "",
				TaskIndex: 1,
			},
		},
	}
	logMap["begindeviceburn"] = time.Now().UnixNano()
	groupid := ""
	for index := 0; index < cd.info.Burn.RetryTimes; index++ {
		groupid, err = api.ADriver.DeviceBurn(parameter, token)
		if err != nil {
			err := fmt.Errorf("api.ADriver.DeviceBurn username {%v} error {%v}", username, err)
			log.Error(err)
		} else {
			break
		}
	}
	if err != nil {
		endinfo.Err = err
	}
	logMap["entertaskswait"] = time.Now().UnixNano()
	logMap["groupid"] = groupid

	// 监听结束
	err = <-endchan
	if err != nil {
		err := fmt.Errorf("endchan username {%v} error {%v}", username, err)
		log.Error(err)
		endinfo.Err = err
		return
	}

	// 日志记录
	logbyte, err := json.MarshalIndent(logMap, "", "		")
	err = ioutil.WriteFile(filepath.Join(datadir, fmt.Sprintf("%s.json", username)), logbyte, 0644)
	if err != nil {
		err := fmt.Errorf("ioutil.WriteFile username {%v} error {%v}", username, err)
		log.Error(err)
		endinfo.Err = err
		return
	}

}
