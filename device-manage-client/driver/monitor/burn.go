package monitor

import (
	"fmt"
	"linklab/device-control-v2/base-library/client/iotnode/api"
	"linklab/device-control-v2/base-library/client/iotnode/topichandler"
	"linklab/device-control-v2/device-manage-client/driver/iotnode"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) burnprocess() {
	// 处理烧写任务的过程中的任何错误,都判定为烧写任务失败
	// 如果无法成功调用云端接口,则直接发送错误信号
	for burninfo := range md.burnchan {

		burntime := topichandler.BurnTime{
			BeginDownloadFile: time.Now(),
			EndDownloadFile:   time.Now(),
			BeginBurn:         time.Now(),
			EndBurn:           time.Now(),
		}

		// 解析开发板类型
		board, err := iotnode.IDriver.GetBoardFromDevPort(burninfo.DeviceID)
		if err != nil {
			err = fmt.Errorf("get board from devport {%v} error {%v}", burninfo.DeviceID, err)
			log.Error(err)
			if err = topichandler.TDriver.PubBurnResult(burninfo, &burntime, false, err.Error()); err != nil {
				log.Errorf("burn result upload error {%v}", err)
				md.errchan <- err
			}
			continue
		}

		burntime.BeginDownloadFile = time.Now()

		// 下载待烧写文件
		// 考虑到流量较多时,存在下载困难的问题,因此下载烧写文件需要多次尝试
		filePath, err := "", nil
		for retrycount := 0; retrycount < md.info.Burn.MaxFileDownloadRetry; retrycount++ {

			filePath, err = api.ADriver.FileDownload(board, burninfo.FileHash, burninfo.GroupID, burninfo.TaskIndex, "hex")
			if err == nil {
				break
			}

			err = fmt.Errorf("burn file {%v:%v} download error {%v}", board, burninfo.FileHash, err)
			log.Error(err)

			time.Sleep(time.Second * time.Duration(md.info.Burn.FileDownloadRetryInterval))
		}

		if err != nil {
			err = fmt.Errorf("burn file download error {%v}, max file download times {%v}", err, md.info.Burn.MaxFileDownloadRetry)
			log.Error(err)
			if err = topichandler.TDriver.PubBurnResult(burninfo, &burntime, false, err.Error()); err != nil {
				log.Errorf("burn result upload error {%v}", err)
				md.errchan <- err
			}
		}

		burntime.EndDownloadFile = time.Now()
		burntime.BeginBurn = time.Now()

		// 开始烧写
		if _, err := iotnode.IDriver.Burn(burninfo, filePath); err != nil {
			err = fmt.Errorf("devport {%v} burn error {%v}", burninfo.DeviceID, err)
			log.Error(err)
			if err = topichandler.TDriver.PubBurnResult(burninfo, &burntime, false, err.Error()); err != nil {
				log.Errorf("burn result upload error {%v}", err)
				md.errchan <- err
			}
			continue
		}

		burntime.EndBurn = time.Now()

		// 上传烧写成功
		if err = topichandler.TDriver.PubBurnResult(burninfo, &burntime, true, "success"); err != nil {
			log.Errorf("burn result upload error {%v}", err)
			md.errchan <- err
		}
	}
}

func (md *Driver) burnstartup() {
	for index := 0; index < runtime.NumCPU()*md.info.Burn.ThreadMultiple; index++ {

		log.Debugf("burn process {%v} start up", index)
		go md.burnprocess()
	}
}
