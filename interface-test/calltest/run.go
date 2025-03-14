package calltest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

// Run 运行接口测试实例
func (cd *Driver) Run(requesttimes int) error {

	// 数据文件目录创建
	datadir := filepath.Join(cd.info.Test.DataDir, time.Now().Format(time.RFC3339Nano))
	if _, err := os.Stat(datadir); os.IsExist(err) {
		// 数据文件目录已经存在,返回错误
		err := fmt.Errorf("datadir {%v} already exists", datadir)
		log.Error(err)
		return err
	}

	if err := os.Mkdir(datadir, 0755); err != nil {
		err := fmt.Errorf("os.Mkdir {%v} error {%v}", datadir, err)
		log.Error(err)
		return err
	}

	// 交换目录创建
	tmpdir := filepath.Join(datadir, "tmp")
	if err := os.RemoveAll(tmpdir); err != nil {
		err := fmt.Errorf("removeall {%v} error {%v}", tmpdir, err)
		log.Error(err)
		return err
	}

	if err := os.Mkdir(tmpdir, 0755); err != nil {
		err := fmt.Errorf("os.Mkdir {%v} error {%v}", datadir, err)
		log.Error(err)
		return err
	}

	// 日志记录开始时间
	logMap := make(map[string]interface{})
	logMap["begintime"] = time.Now().Format(time.RFC3339Nano)
	logMap["begintime.utc"] = time.Now().UTC().Format(time.RFC3339)

	// 启动所有测试单例
	for count := 0; count < requesttimes; count++ {
		go cd.testthread(datadir, tmpdir, count)
	}

	// 等待所有测试单例完成
	for count := 0; count < requesttimes; count++ {
		select {
		case endinfo := <-cd.endchan:
			if endinfo == nil {
				err := errors.New("endinfo nil error")
				log.Error(err)
				return err
			}
			if endinfo.Err != nil {
				err := fmt.Errorf("endinfo username {%v} error {%v}", endinfo.UserName, endinfo.Err)
				log.Error(err)
				return err
			}
			log.Infof("endinfo username {%v} success", endinfo.UserName)
		}
	}

	// 解析所有测试单例日志
	// avglogbytes := 0.0
	alllogbytes := 0
	avguserloginresp := time.Duration(0)
	avgfileuploadresp := time.Duration(0)
	avgentertaskswaitresp := time.Duration(0)
	avgdeviceallocateresp := time.Duration(0)
	avgdeviceburnresp := time.Duration(0)
	avgdevicerunresp := time.Duration(0)
	avgcompileresp := time.Duration(0)
	alllogrecvcount := 0
	alllogrecvresp := 0
	alldevicelogrecvresp := 0
	allnatslogrecvresp := 0
	alluserlogrecvresp := 0

	lostallocatecount := 0

	for index := 0; index < requesttimes; index++ {
		slogbyte, err := ioutil.ReadFile(fmt.Sprintf("%s/%s-%v.json", datadir, cd.info.Login.UserName, index))
		if err != nil {
			err := fmt.Errorf("ioutil.ReadFile error {%v}", err)
			log.Error(err)
			return err
		}
		tasktestresult := TaskTestResult{}
		if err := json.Unmarshal(slogbyte, &tasktestresult); err != nil {
			err := fmt.Errorf("json.Unmarshal error {%v}", err)
			log.Error(err)
			return err
		}

		alllogbytes += tasktestresult.LogBytesCount

		beginuserlogin := time.Unix(0, tasktestresult.BeginUserLogin)
		enduserlogin := time.Unix(0, tasktestresult.EndUserLogin)
		avguserloginresp += enduserlogin.Sub(beginuserlogin)

		beginfileupload := time.Unix(0, tasktestresult.BeginFileUpload)
		endfileupload := time.Unix(0, tasktestresult.EndFileUpload)
		avgfileuploadresp += endfileupload.Sub(beginfileupload)

		begindeviceburn := time.Unix(0, tasktestresult.BeginDeviceBurn)
		entertaskswait := time.Unix(0, tasktestresult.EnterTasksWait)

		if entertaskswait.Sub(begindeviceburn) > 0 {
			avgentertaskswaitresp += entertaskswait.Sub(begindeviceburn)
		}

		// log.Debugf("avgentertaskswaitresp{%v} avg{%v} single{%v}", index, avgentertaskswaitresp, entertaskswait.Sub(begindeviceburn))

		enddeviceburn := time.Unix(0, tasktestresult.EndDeviceBurn)
		if tasktestresult.DeviceAllocate < 1 {
			lostallocatecount++
		} else {

			deviceallocate := time.Unix(0, tasktestresult.DeviceAllocate)
			if deviceallocate.Sub(begindeviceburn) > 0 {
				avgdeviceallocateresp += deviceallocate.Sub(begindeviceburn)
				log.Debugf("avgdeviceallocateresp{%v} avg{%v} single{%v}", index, avgdeviceallocateresp, deviceallocate.Sub(begindeviceburn))
			}

			avgdeviceburnresp += enddeviceburn.Sub(deviceallocate)
			log.Debugf("avgdeviceburnresp{%v} avg{%v} single{%v}", index, avgdeviceburnresp, enddeviceburn.Sub(deviceallocate))
		}

		enddevicerun := time.Unix(0, tasktestresult.EndDeviceRun)
		avgdevicerunresp += enddevicerun.Sub(enddeviceburn)

		logrecvcount := tasktestresult.LogRecvCount
		alllogrecvcount += logrecvcount

		peravglogrecvresp := tasktestresult.AvgLogRecvResp
		alllogrecvresp += (peravglogrecvresp * logrecvcount)

		peravgdevicelogrecvresp := tasktestresult.AvgDeviceLogRecvResp
		alldevicelogrecvresp += (peravgdevicelogrecvresp * logrecvcount)

		allnatslogrecvresp += (tasktestresult.AvgNatsLogRecvResp * logrecvcount)
		alluserlogrecvresp += (tasktestresult.AvgUserLogRecvResp * logrecvcount)

		begincompile := time.Unix(0, tasktestresult.BeginCompile)
		endcompile := time.Unix(0, tasktestresult.EndCompile)
		avgcompileresp += endcompile.Sub(begincompile)
	}
	avglogbytes := alllogbytes / requesttimes
	avguserloginresp /= time.Duration(requesttimes)
	avgfileuploadresp /= time.Duration(requesttimes)
	avgentertaskswaitresp /= time.Duration(requesttimes)
	avgdeviceallocateresp /= time.Duration(requesttimes - lostallocatecount)
	avgdeviceburnresp /= time.Duration(requesttimes - lostallocatecount)
	avgdevicerunresp /= time.Duration(requesttimes)
	avgcompileresp /= time.Duration(requesttimes)
	logMap["avglogbytes"] = avglogbytes
	logMap["avguserloginresp"] = int64(avguserloginresp / time.Millisecond)
	logMap["avgfileuploadresp"] = int64(avgfileuploadresp / time.Millisecond)
	logMap["avgentertaskswaitresp"] = int64(avgentertaskswaitresp / time.Millisecond)
	logMap["avgdeviceallocateresp"] = int64(avgdeviceallocateresp / time.Millisecond)
	logMap["avgdeviceburnresp"] = int64(avgdeviceburnresp / time.Millisecond)
	logMap["avgdevicerunresp"] = int64(avgdevicerunresp / time.Millisecond)
	logMap["avglogrecvcount"] = float64(alllogrecvcount) / float64(requesttimes)
	logMap["avglogrecvresp"] = alllogrecvresp / alllogrecvcount
	logMap["avgdevicelogrecvresp"] = alldevicelogrecvresp / alllogrecvcount
	logMap["avgnatslogrecvresp"] = allnatslogrecvresp / alllogrecvcount
	logMap["avguserlogrecvresp"] = alluserlogrecvresp / alllogrecvcount
	logMap["lostallocatecount"] = lostallocatecount
	logMap["avgcompileresp"] = int64(avgcompileresp / time.Millisecond)

	// 日志记录结束时间
	logMap["endtime"] = time.Now().Format(time.RFC3339Nano)
	logMap["endtime.utc"] = time.Now().UTC().Format(time.RFC3339)

	logbyte, err := json.MarshalIndent(logMap, "", "		")
	err = ioutil.WriteFile(filepath.Join(datadir, "info.json"), logbyte, 0644)
	if err != nil {
		err := fmt.Errorf("ioutil.WriteFile error {%v}", err)
		log.Error(err)
		return err
	}

	// 导出数据库日志
	time.Sleep(time.Minute)
	cmd := exec.Command("sh", "-c", cd.info.OSSQuery.PodMetrics)
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("DATA_DIR=%s", datadir),
		fmt.Sprintf("START_TIME=%s", logMap["begintime.utc"]),
		fmt.Sprintf("END_TIME=%s", logMap["endtime.utc"]),
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		err := fmt.Errorf("cmd.CombinedOutput {%v} error {%v}", string(stdoutStderr), err)
		log.Error(err)
		return err
	}

	cmd = exec.Command("sh", "-c", cd.info.OSSQuery.NodeMetrics)
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("DATA_DIR=%s", datadir),
		fmt.Sprintf("START_TIME=%s", logMap["begintime.utc"]),
		fmt.Sprintf("END_TIME=%s", logMap["endtime.utc"]),
	)
	stdoutStderr, err = cmd.CombinedOutput()
	if err != nil {
		err := fmt.Errorf("cmd.CombinedOutput {%v} error {%v}", string(stdoutStderr), err)
		log.Error(err)
		return err
	}

	return nil
}
