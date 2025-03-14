package aiboxnode

import (
	"bufio"
	"fmt"
	"io"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

// Exec 执行
func (ad *Driver) Exec(burninfo *msg.ClientBurnMsg, execzip string) (int, error) {

	// 解析开发板类型
	board, err := ad.GetBoardFromDevPort(burninfo.DeviceID)
	if err != nil {
		err := fmt.Errorf("get board from devport {%v} error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return 0, err
	}

	// 设备操作加锁
	value, isOk := ad.devices.devicesMap.Load(burninfo.DeviceID)
	if isOk == false {
		err := fmt.Errorf("devport {%v} not in devicesMap error", burninfo.DeviceID)
		log.Error(err)
		return 0, err
	}
	devicestatus := value.(*DeviceStatus)
	devicestatus.Lock.Lock()
	defer func() {
		devicestatus.Lock.Unlock()
	}()

	// 设置正在执行状态
	if err := ad.setrun(burninfo.DeviceID); err != nil {
		err = fmt.Errorf("devport {%v} set run state error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return 0, err
	}
	defer func() {
		if err := ad.setidle(burninfo.DeviceID); err != nil {
			err = fmt.Errorf("devport {%v} set idle state error {%v}", burninfo.DeviceID, err)
			log.Error(err)
		}
	}()
	// 设置执行信息
	devicestatus.BurnInfo = burninfo

	// 确定工作目录路径
	workdir := filepath.Join(ad.info.Workspace, fmt.Sprintf("%s-%v", burninfo.GroupID, burninfo.TaskIndex))

	// 解压zip文件
	if err := archiver.Unarchive(execzip, workdir); err != nil {
		err = fmt.Errorf("archiver.Unarchive {%v} {%v} error {%v}", execzip, workdir, err)
		log.Error(err)
		return 0, err
	}

	// 运行结束时需要将文件删除
	defer func() {
		if err := os.Remove(execzip); err != nil {
			log.Errorf("execzip {%v} remove error {%v}", execzip, err)
		}
		if err := os.RemoveAll(workdir); err != nil {
			log.Errorf("workdir {%v} remove error {%v}", workdir, err)
		}
	}()

	// 获取日志写入通道
	logchan, err := ad.getDeviceWriteLogChan(burninfo.DeviceID)
	if err != nil {
		err = fmt.Errorf("ad.getDeviceWriteLogChan {%v} error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return 0, err
	}

	log.Debugf("exec prepare")

	cmdargs := strings.Fields(ad.info.Boards[board].ExecCmd)
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
	cmd.Dir = workdir
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("cmd.StdoutPipe error {%v}", err)
		log.Error(err)
		return 0, err
	}
	linereader := bufio.NewReader(stdout)

	time.Sleep(time.Second)
	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("cmd.Start error {%v}", err)
		log.Error(err)
		return 0, err
	}

	log.Debugf("exec start")

	timeoutflag := false
	donechan := make(chan error)
	go func() {
		for {

			if timeoutflag {
				log.Debugf("timeoutflag true")
				return
			}

			line, err := linereader.ReadString('\n')
			if err != nil {
				log.Errorf("linereader.ReadString error {%v}", err)
				if err == io.EOF {
					donechan <- nil
					return
				}
				donechan <- err
				return
			}

			select {
			case logchan <- &LogMsg{
				Msg:       line,
				TimeStamp: time.Now().UnixNano(),
			}:
			case <-time.After(time.Millisecond * time.Duration(ad.info.DeviceLog.LogTimeOutMill)):
				log.Errorf("devport {%v} log {%v} into logchan timeout {%vmills} error", burninfo.DeviceID, line, ad.info.DeviceLog.LogTimeOutMill)
				continue
			}
		}
	}()

	defer func() {
		log.Debugf("exec end")
	}()

	select {
	case err := <-donechan:
		if err == nil {
			return 0, nil
		}
		err = fmt.Errorf("process run error {%v}", err)
		log.Error(err)
		return 0, err
	case <-time.After(time.Second * time.Duration(burninfo.RunTime)):
		timeoutflag = true
		if err := cmd.Process.Kill(); err != nil {
			err = fmt.Errorf("process timeout kill error {%v}", err)
			log.Error(err)
			return 1, err
		}
		log.Debugf("cmd.Process.Kill success")
		return 1, nil
	}
}
