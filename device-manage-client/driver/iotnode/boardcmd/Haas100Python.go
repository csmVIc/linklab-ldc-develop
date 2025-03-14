package boardcmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/tool"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/albenik/go-serial/v2"
	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

// Haas100Python 操作
type Haas100Python struct {
	baudrate int
	ssid     string
	password string
}

// Burn 烧写
func (driver *Haas100Python) Burn(devport string, burnfile string, burncmd string, serialport *serial.Port) (string, error) {

	// 判断串口是否打开,如果打开,则将串口关闭
	if err := driver.CloseSerial(serialport); err != nil {
		err := fmt.Errorf("close serial error {%v}", err)
		log.Error(err)
		return "", err
	}

	// zip
	zipfile := burnfile + ".zip"
	zipdata, err := ioutil.ReadFile(burnfile)
	if err != nil {
		err = fmt.Errorf("ioutil.ReadFile {%v} error {%v}", burnfile, err)
		log.Error(err)
		return "", err
	}
	err = ioutil.WriteFile(zipfile, zipdata, 0644)
	if err != nil {
		err = fmt.Errorf("ioutil.WriteFile {%v} error {%v}", zipfile, err)
		log.Error(err)
		return "", err
	}
	defer func() {
		if err := os.Remove(zipfile); err != nil {
			log.Errorf("os.Remove {%v} remove error", zipfile)
		}
	}()

	// 产生提取目录
	extractdir := filepath.Join("./tmp", tool.GenerateRandomString())
	defer func() {
		if err := os.RemoveAll(extractdir); err != nil {
			log.Errorf("extractdir {%v} remove error", extractdir)
		}
	}()

	// 解压目录
	if err := archiver.Unarchive(zipfile, extractdir); err != nil {
		err = fmt.Errorf("archiver.Unarchive {%v} {%v} error {%v}", zipfile, extractdir, err)
		log.Error(err)
		return "", err
	}

	// 重启
	if err := driver.Reset(devport, serialport); err != nil {
		err = fmt.Errorf("Reset {%v} error {%v}", devport, err)
		log.Error(err)
		return "", err
	}

	// 获取ip地址
	ipv4str, err := driver.getIpv4(devport, serialport)
	log.Debugf("get ipv4str {%v}", ipv4str)
	if err != nil {
		err = fmt.Errorf("getIpv4 {%v} error {%v}", devport, err)
		log.Error(err)
		return "", err
	}

	// 启动Tftp server
	err = driver.startTftpServer(devport, serialport)
	if err != nil {
		err = fmt.Errorf("startTftpServer {%v} error {%v}", devport, err)
		log.Error(err)
		return "", err
	}

	// 上传文件
	datafiles, err := ioutil.ReadDir(extractdir)
	if err != nil {
		err = fmt.Errorf("ioutil.ReadDir {%v} error {%v}", extractdir, err)
		log.Error(err)
		return "", err
	}

	pyfilename := ""
	for _, datafile := range datafiles {
		if strings.HasSuffix(datafile.Name(), ".py") || strings.HasSuffix(datafile.Name(), ".png") || strings.HasSuffix(datafile.Name(), ".bmp") {

			dstfilename := "test"
			if strings.HasSuffix(datafile.Name(), ".py") {
				dstfilename += ".py"
				pyfilename = "test.py"
			} else if strings.HasSuffix(datafile.Name(), ".png") {
				dstfilename += ".png"
			} else if strings.HasSuffix(datafile.Name(), ".bmp") {
				dstfilename += ".bmp"
			}

			datafilepath := filepath.Join(extractdir, datafile.Name())
			cmdstr := fmt.Sprintf(burncmd, ipv4str, datafilepath, dstfilename)
			log.Infof("Haas100Python send file cmd {%v}", cmdstr)
			cmd := exec.Command("sh", "-c", cmdstr)
			stdoutStderr, err := cmd.CombinedOutput()
			log.Debugf("tftp shell {%v}", string(stdoutStderr))
			if err != nil {
				err = fmt.Errorf("cmd.CombinedOutput cmd {%v} {%v} error {%v}", cmdstr, string(stdoutStderr), err)
				log.Error(err)
				return "", err
			}

			if strings.Index(string(stdoutStderr), "Error code 2") >= 0 {
				err = fmt.Errorf("cmd.CombinedOutput cmd {%v} {%v} error {%v}", cmdstr, string(stdoutStderr), err)
				log.Error(err)
				return "", err
			}

			time.Sleep(time.Millisecond * 100)
		}
	}

	// log.Debug("pyfilename ", pyfilename)
	if len(pyfilename) >= 0 {
		err = driver.pythonRun(devport, serialport, pyfilename)
		if err != nil {
			err := fmt.Errorf("pythonRun {%v} error {%v}", devport, err)
			log.Error(err)
			return "", err
		}
		return "", nil
	}

	err = errors.New("can't find python file error")
	log.Error(err)
	return "", err
}

// startTftpServer 启动tftp服务器
func (driver *Haas100Python) startTftpServer(devport string, serialport *serial.Port) error {

	// 判断串口是否打开,如果没有打开,则重新打开串口
	if err := driver.ReOpen(serialport); err != nil {
		err := fmt.Errorf("serialport {%v} reopen error {%v}", devport, err)
		log.Error(err)
		return err
	}

	// 写入启动命令
	cmdstr := "tftp server start\r\n"
	log.Debugf("python cmd {%v}", cmdstr)
	if _, err := serialport.Write([]byte(cmdstr)); err != nil {
		err := fmt.Errorf("serialport {%v} write {%v} error {%v}", devport, cmdstr, err)
		log.Error(err)
		return err
	}

	// 读取tftp结果
	isfind := -1
	logbuffer := ""
	maxwaittime := time.Now().Add(time.Second * 30)
	for isfind < 0 {
		readready, err := serialport.ReadyToRead()
		if err != nil {
			err := fmt.Errorf("serialport {%v} ReadyToRead error {%v}", devport, err)
			log.Error(err)
			return err
		}

		if readready > 0 {
			tmpbuffer := make([]byte, 1024)
			readlen, err := serialport.Read(tmpbuffer)
			if err != nil {
				err := fmt.Errorf("serialport {%v} Read error {%v}", devport, err)
				log.Error(err)
				return err
			}
			logbuffer = logbuffer + string(tmpbuffer[:readlen])

			lines := strings.Split(logbuffer, "\n")
			for i := 0; i < len(lines); i++ {
				log.Debugf("tftp server info {%v}", lines[i])
				isfind = strings.Index(lines[i], "tftp start server done")
				if isfind >= 0 {
					return nil
				}
			}
			logbuffer = lines[len(lines)-1]
		}

		// 等待tftp server启动超时
		if maxwaittime.Before(time.Now()) {
			err := fmt.Errorf("wait tftp start server error")
			log.Error(err)
			return err
		}
	}

	return nil
}

// getIpv4 获取ip地址
func (driver *Haas100Python) getIpv4(devport string, serialport *serial.Port) (string, error) {

	// 判断串口是否打开,如果没有打开,则重新打开串口
	if err := driver.ReOpen(serialport); err != nil {
		err := fmt.Errorf("serialport {%v} reopen error {%v}", devport, err)
		log.Error(err)
		return "", err
	}

	isfind := -1
	logbuffer := ""
	maxwaittime := time.Now().Add(time.Second * 30)
	for isfind < 0 {
		readready, err := serialport.ReadyToRead()
		if err != nil {
			err := fmt.Errorf("serialport {%v} ReadyToRead error {%v}", devport, err)
			log.Error(err)
			return "", err
		}

		if readready > 0 {
			tmpbuffer := make([]byte, 1024)
			readlen, err := serialport.Read(tmpbuffer)
			if err != nil {
				err := fmt.Errorf("serialport {%v} Read error {%v}", devport, err)
				log.Error(err)
				return "", err
			}
			logbuffer = logbuffer + string(tmpbuffer[:readlen])

			lines := strings.Split(logbuffer, "\n")
			for i := 0; i < len(lines)-1; i++ {
				log.Debugf("getIpv4 info {%v}", lines[i])
				isfind = strings.Index(lines[i], "micropython")
				if isfind >= 0 {
					break
				}
			}
			logbuffer = lines[len(lines)-1]
		}

		// 等待micropython初始超时错误
		if maxwaittime.Before(time.Now()) {
			err := fmt.Errorf("wait micropython initialization timeout error")
			log.Error(err)
			return "", err
		}
	}

	// 写入联网命令
	cmdstr := fmt.Sprintf("\r\npython /data/python/network/ConnectWifi.py %s %s\r\n", driver.ssid, driver.password)
	log.Debugf("python cmd {%v}", cmdstr)
	if _, err := serialport.Write([]byte(cmdstr)); err != nil {
		err := fmt.Errorf("serialport {%v} write {%v} error {%v}", devport, cmdstr, err)
		log.Error(err)
		return "", err
	}

	ipv4str := ""
	logbuffer = ""
	maxwaittime = time.Now().Add(time.Second * 30)
	for len(ipv4str) < 1 {
		readready, err := serialport.ReadyToRead()
		if err != nil {
			err := fmt.Errorf("serialport {%v} ReadyToRead error {%v}", devport, err)
			log.Error(err)
			return "", err
		}

		if readready > 0 {
			tmpbuffer := make([]byte, 1024)
			readlen, err := serialport.Read(tmpbuffer)
			if err != nil {
				err := fmt.Errorf("serialport {%v} Read error {%v}", devport, err)
				log.Error(err)
				return "", err
			}

			// log.Debugln(string(tmpbuffer[:readlen]))
			logbuffer = logbuffer + string(tmpbuffer[:readlen])
			lines := strings.Split(logbuffer, "\n")
			for i := 0; i < len(lines)-1; i++ {
				log.Debugf("burn log info {%v}", lines[i])
				ipfindindex := strings.Index(lines[i], "ip_address=")
				if ipfindindex >= 0 {
					ipv4str = lines[i][len("ip_address=") : len(lines[i])-1]
					log.Debugf("ipv4str find {%v}", ipv4str)
					break
				}
			}
			logbuffer = lines[len(lines)-1]
		}

		// 等待联网IP地址错误
		if maxwaittime.Before(time.Now()) {
			err := fmt.Errorf("wait ip_address timeout error")
			log.Error(err)
			return "", err
		}
	}

	return ipv4str, nil
}

// pythonRun 运行python脚本
func (driver *Haas100Python) pythonRun(devport string, serialport *serial.Port, pyfile string) error {

	// 判断串口是否打开,如果没有打开,则重新打开串口
	if err := driver.ReOpen(serialport); err != nil {
		err := fmt.Errorf("serialport {%v} reopen error {%v}", devport, err)
		log.Error(err)
		return err
	}

	cmdstr := fmt.Sprintf("python /data/%v\r\n", pyfile)
	log.Debugf("python cmd {%v}", cmdstr)
	if _, err := serialport.Write([]byte(cmdstr)); err != nil {
		err := fmt.Errorf("serialport {%v} write {%v} error {%v}", devport, cmdstr, err)
		log.Error(err)
		return err
	}

	return nil
}

// Reset 重置
func (driver *Haas100Python) Reset(devport string, serialport *serial.Port) error {

	// 判断串口是否打开,如果没有打开,则重新打开串口
	if err := driver.ReOpen(serialport); err != nil {
		err := fmt.Errorf("serialport {%v} reopen error {%v}", devport, err)
		log.Error(err)
		return err
	}

	if _, err := serialport.Write([]byte("reboot\r\n")); err != nil {
		err := fmt.Errorf("serialport {%v} write reboot error {%v}", devport, err)
		log.Error(err)
		return err
	}

	// 清空读取buffer
	if err := serialport.ResetInputBuffer(); err != nil {
		err := fmt.Errorf("Haas100Python reset input buffer error {%v}", err)
		log.Error(err)
		return err
	}

	// 关闭串口
	if err := driver.CloseSerial(serialport); err != nil {
		err := fmt.Errorf("Haas100Python CloseSerial error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}

// GetBaudRate 获取波特率
func (driver *Haas100Python) GetBaudRate() int {
	return driver.baudrate
}

// SetBaudRate 设置波特率
func (driver *Haas100Python) SetBaudRate(baudrate int) {
	driver.baudrate = baudrate
}

// OpenSerial 打开串口
func (driver *Haas100Python) OpenSerial(devport string) (*serial.Port, error) {
	serialport, err := serial.Open(devport, serial.WithBaudrate(driver.baudrate), serial.WithParity(serial.NoParity), serial.WithDataBits(8), serial.WithStopBits(serial.OneStopBit))
	if err != nil {
		err = fmt.Errorf("devport {%v} open error {%v}", devport, err)
		log.Error(err)
		return nil, err
	}
	return serialport, nil
}

// CloseSerial 关闭串口
func (driver *Haas100Python) CloseSerial(serialport *serial.Port) error {
	if serialport == nil {
		err := errors.New("close serialport nil error")
		log.Error(err)
		return err
	}
	if driver.IsOpen(serialport) == false {
		return nil
	}
	if err := serialport.Close(); err != nil {
		err = fmt.Errorf("serialport.Close error {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}

// IsOpen 判断串口是否打开
func (driver *Haas100Python) IsOpen(serialport *serial.Port) bool {
	if serialport == nil {
		return false
	}
	_, err := serialport.ReadyToRead()
	if err != nil {
		return false
	}
	return true
}

// ReOpen 重新
func (driver *Haas100Python) ReOpen(serialport *serial.Port) error {
	if serialport == nil {
		err := errors.New("reopen serialport nil error")
		log.Error(err)
		return err
	}
	if driver.IsOpen(serialport) == true {
		return nil
	}
	tmpptr, err := driver.OpenSerial(serialport.String())
	if err != nil {
		err := fmt.Errorf("open serial {%v} error {%v}", serialport.String(), err)
		log.Error(err)
		return err
	}
	*serialport = *tmpptr
	return nil
}

// SetResetCmd 设置重置命令
func (driver *Haas100Python) SetResetCmd(resetcmd string) {
	return
}

// BurnEmptyProgram 烧写空程序
func (driver *Haas100Python) BurnEmptyProgram(devport string, burncmd string, serialport *serial.Port) error {
	if err := driver.Reset(devport, serialport); err != nil {
		err = fmt.Errorf("reset error {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}

// SetEmptyProgram 设置空程序
func (driver *Haas100Python) SetEmptyProgram(emptyprogram string) {
	return
}

// SetWifi 设置WiFi
func (driver *Haas100Python) SetWifi(ssid string, password string) {
	driver.ssid = ssid
	driver.password = password
}
