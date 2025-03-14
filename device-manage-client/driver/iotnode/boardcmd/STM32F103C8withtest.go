package boardcmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/albenik/go-serial/v2"
	log "github.com/sirupsen/logrus"
)

// STM32F103C8with4G 操作
type STM32F103C8with4G struct {
	baudrate     int
	emprtprogram string
}

// Burn 烧写
func (driver *STM32F103C8with4G) Burn(devport string, burnfile string, burncmd string, serialport *serial.Port) (string, error) {

	// 判断串口是否打开,如果打开,则将串口关闭
	if err := driver.CloseSerial(serialport); err != nil {
		err := fmt.Errorf("close serial error {%v}", err)
		log.Error(err)
		return "", err
	}

	cmdstr := fmt.Sprintf(burncmd, burnfile, devport)
	log.Infof("STM32F103C8with4G burn cmd {%v}", cmdstr)

	cmd := exec.Command("sh", "-c", cmdstr)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Infof("STM32F103C8with4G burn result {%v}", string(stdoutStderr))
	if err != nil {
		log.Errorf("STM32F103C8with4G burn error {%v}", err)
		return fmt.Sprintf("%s\n\n%s", cmdstr, string(stdoutStderr)), err
	}

	// 增加判断错误
	if sindex := strings.Index(string(stdoutStderr), "success"); sindex < 0 {
		err = fmt.Errorf("STM32F103C8with4G burn error {%v}", string(stdoutStderr))
		return fmt.Sprintf("%s\n\n%s", cmdstr, string(stdoutStderr)), err
	}

	return fmt.Sprintf("%s\n\n%s", cmdstr, string(stdoutStderr)), nil
}

// Reset 重置
func (driver *STM32F103C8with4G) Reset(devport string, serialport *serial.Port) error {

	// 判断串口是否打开,如果没有打开,则重新打开串口
	if err := driver.ReOpen(serialport); err != nil {
		err := fmt.Errorf("serialport {%v} reopen error {%v}", devport, err)
		log.Error(err)
		return err
	}

	// 重置
	time.Sleep(time.Millisecond * 100)
	if err := serialport.SetDTR(false); err != nil {
		err := fmt.Errorf("STM32F103C8with4G set dtr error {%v}", err)
		log.Error(err)
		return err
	}
	time.Sleep(time.Millisecond * 200)
	if err := serialport.SetRTS(false); err != nil {
		err := fmt.Errorf("STM32F103C8with4G set dtr error {%v}", err)
		log.Error(err)
		return err
	}

	// 清空读取buffer
	if err := serialport.ResetInputBuffer(); err != nil {
		err := fmt.Errorf("STM32F103C8with4G reset input buffer error {%v}", err)
		log.Error(err)
		return err
	}

	// 关闭串口
	if err := driver.CloseSerial(serialport); err != nil {
		err := fmt.Errorf("STM32F103C8with4G CloseSerial error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}

// GetBaudRate 获取波特率
func (driver *STM32F103C8with4G) GetBaudRate() int {
	return driver.baudrate
}

// SetBaudRate 设置波特率
func (driver *STM32F103C8with4G) SetBaudRate(baudrate int) {
	driver.baudrate = baudrate
}

// OpenSerial 打开串口
func (driver *STM32F103C8with4G) OpenSerial(devport string) (*serial.Port, error) {
	serialport, err := serial.Open(devport, serial.WithBaudrate(driver.baudrate), serial.WithParity(serial.NoParity), serial.WithDataBits(8), serial.WithStopBits(serial.OneStopBit))
	if err != nil {
		err = fmt.Errorf("devport {%v} open error {%v}", devport, err)
		log.Error(err)
		return nil, err
	}
	return serialport, nil
}

// CloseSerial 关闭串口
func (driver *STM32F103C8with4G) CloseSerial(serialport *serial.Port) error {
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
func (driver *STM32F103C8with4G) IsOpen(serialport *serial.Port) bool {
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
func (driver *STM32F103C8with4G) ReOpen(serialport *serial.Port) error {
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
func (driver *STM32F103C8with4G) SetResetCmd(resetcmd string) {
	return
}

// BurnEmptyProgram 烧写空程序
func (driver *STM32F103C8with4G) BurnEmptyProgram(devport string, burncmd string, serialport *serial.Port) error {
	if output, err := driver.Burn(devport, driver.emprtprogram, burncmd, serialport); err != nil {
		err = fmt.Errorf("burn error {%v} {%v}", output, err)
		log.Error(err)
		return err
	}
	return nil
}

// SetEmptyProgram 设置空程序
func (driver *STM32F103C8with4G) SetEmptyProgram(emptyprogram string) {
	driver.emprtprogram = emptyprogram
}

// SetWifi 设置WiFi
func (driver *STM32F103C8with4G) SetWifi(ssid string, password string) {
}
