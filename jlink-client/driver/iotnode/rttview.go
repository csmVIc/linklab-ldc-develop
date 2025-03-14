package iotnode

import (
	"bufio"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/jlink-client/driver/iotnode/boardcmd"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// RTTView 日志获取和命令交互
func (id *Driver) RTTView(burninfo *msg.ClientBurnMsg) (int, error) {

	// 解析开发板类型
	board, err := id.GetBoardFromDevPort(burninfo.DeviceID)
	if err != nil {
		err := fmt.Errorf("get board from devport {%v} error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return 0, err
	}

	// 设备操作加锁
	value, isOk := id.devices.devicesMap.Load(burninfo.DeviceID)
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

	// 退出设置设备为空闲状态
	defer func() {
		if err := id.setidle(burninfo.DeviceID); err != nil {
			err = fmt.Errorf("devport {%v} set idle state error {%v}", burninfo.DeviceID, err)
			log.Error(err)
		}
	}()

	// 退出时需要烧写空白程序
	defer func() {
		if err := id.boardCmdMap[board].BurnEmptyProgram(burninfo.DeviceID, id.info.Commands[board].Burn); err != nil {
			err = fmt.Errorf("devport {%v} burn empty program error {%v}", burninfo.DeviceID, err)
			log.Error(err)
		}
	}()

	// // 获取日志写入通道
	logchan, err := id.getDeviceWriteLogChan(burninfo.DeviceID)
	if err != nil {
		err = fmt.Errorf("id.getDeviceWriteLogChan {%v} error {%v}", burninfo.DeviceID, err)
		log.Error(err)
		return 0, err
	}

	log.Debugf("rtt view prepare")

	serialnum, err := boardcmd.GetJLinkSerialNumber(burninfo.DeviceID)
	if err != nil {
		err = fmt.Errorf("boardcmd.GetJLinkSerialNumber error {%v}", err)
		log.Error(err)
		return 0, err
	}

	cmdstr := fmt.Sprintf(id.info.Commands[board].RTTCmd, serialnum)
	log.Debugln(cmdstr)
	cmdargs := strings.Fields(cmdstr)
	log.Debugln(cmdargs)
	cmd := exec.Command(cmdargs[0], cmdargs[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Dir = "/app"
	cmd.Stderr = cmd.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		err = fmt.Errorf("cmd.StdoutPipe error {%v}", err)
		log.Error(err)
		return 0, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		err = fmt.Errorf("cmd.StdinPipe error {%v}", err)
		log.Error(err)
		return 0, err
	}
	linereader := bufio.NewScanner(stdout)

	err = cmd.Start()
	if err != nil {
		err = fmt.Errorf("cmd.Start error {%v}", err)
		log.Error(err)
		return 0, err
	}

	log.Debugf("rtt view start")

	donechan := make(chan error)
	go func() {
		for linereader.Scan() {
			line := linereader.Text()
			select {
			case logchan <- &LogMsg{
				Msg:       line,
				TimeStamp: time.Now().UnixNano(),
			}:
			case <-time.After(time.Millisecond * time.Duration(id.info.DeviceLog.LogTimeOutMill)):
				log.Errorf("devport {%v} log {%v} into logchan timeout {%vmills} error", burninfo.DeviceID, line, id.info.DeviceLog.LogTimeOutMill)
				continue
			}
		}
		donechan <- nil
		log.Debugf("rtt view scan end")
	}()

	readcontinue := true
	go func() {
		for readcontinue {
			select {
			case cmdstr := <-*devicestatus.CmdChan:
				if _, err := stdin.Write([]byte(cmdstr.Cmd)); err != nil {
					log.Errorf("devport {%v} write cmd {%v} error", burninfo.DeviceID, cmdstr.Cmd)
				}
				log.Debugf("WRITE CMD [%v]", cmdstr.Cmd)
			case <-time.After(time.Second):
				continue
			}
		}
		log.Debugf("rtt view cmd end")
	}()

	defer func() {
		readcontinue = false
		log.Debugf("rtt view exit")
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
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			err = fmt.Errorf("Process.Signal error {%v}", err)
			log.Error(err)
			return 1, err
		}
		time.Sleep(time.Second)
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			err = fmt.Errorf("syscall.Kill error {%v}", err)
			log.Error(err)
			return 1, err
		}
		log.Debugf("cmd.Process.Kill success")
		return 1, nil
	}
}
