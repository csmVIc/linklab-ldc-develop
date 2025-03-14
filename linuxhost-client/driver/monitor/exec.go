package monitor

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/linuxhost-client/driver/linuxhostnode"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) execprocess() {

	defer func() {
		md.endrun = true
		log.Error("exec func set endrun")
	}()

	// 处理烧写任务的过程中的任何错误,都判定为烧写任务失败
	// 如果无法成功调用云端接口,则直接发送错误信号
	for burninfo := range md.burnchan {

		// 解析开发板类型
		board, err := linuxhostnode.LDriver.GetBoardFromDevPort(burninfo.DeviceID)
		if err != nil {
			err = fmt.Errorf("get board from devport {%v} error {%v}", burninfo.DeviceID, err)
			log.Error(err)
			if err = topichandler.TDriver.PubExecErr(burninfo.GroupID, burninfo.TaskIndex, err.Error()); err != nil {
				log.Errorf("exec err upload error {%v}", err)
				md.errchan <- err
			}
			continue
		}

		// 下载待烧写文件
		// 考虑到流量较多时,存在下载困难的问题,因此下载烧写文件需要多次尝试
		filePath, err := "", nil
		for retrycount := 0; retrycount < md.info.Exec.MaxFileDownloadRetry; retrycount++ {
			filePath, err = api.ADriver.FileDownload(board, burninfo.FileHash, burninfo.GroupID, burninfo.TaskIndex, "zip")
			if err == nil {
				break
			}
			err = fmt.Errorf("exec file {%v:%v} download error {%v}", board, burninfo.FileHash, err)
			log.Error(err)
			time.Sleep(time.Second * time.Duration(md.info.Exec.FileDownloadRetryInterval))
		}
		if err != nil {
			err = fmt.Errorf("exec file download error {%v}, max file download times {%v}", err, md.info.Exec.MaxFileDownloadRetry)
			log.Error(err)
			if err = topichandler.TDriver.PubExecErr(burninfo.GroupID, burninfo.TaskIndex, err.Error()); err != nil {
				log.Errorf("exec err upload error {%v}", err)
				md.errchan <- err
			}
		}

		// 开始执行
		istimeout, err := linuxhostnode.LDriver.Exec(burninfo, filePath)
		if err != nil {
			err = fmt.Errorf("process exec error {%v}", err)
			log.Error(err)
			topichandler.TDriver.PubExecErr(burninfo.GroupID, burninfo.TaskIndex, err.Error())
		}

		err = topichandler.TDriver.PubEndRun(burninfo, istimeout, time.Now())
		if err != nil {
			log.Errorf("topichandler.TDriver.PubEndRun error {%v}", err)
			return
		}
	}
}

func (md *Driver) execstartup() {
	for index := 0; index < runtime.NumCPU()*md.info.Exec.ThreadMultiple; index++ {
		log.Debugf("exec process {%v} start up", index)
		go md.execprocess()
	}
}
