package boardcmd

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/albenik/go-serial/v2"
	log "github.com/sirupsen/logrus"
)

// ArduinoUno 操作
type ArduinoUno struct {
	baudrate     int
	emprtprogram string
}

// Burn 烧写
func (driver *ArduinoUno) Burn(devport string, burnfile string, burncmd string, serialport *serial.Port) (string, error) {

	// 判断串口是否打开,如果打开,则将串口关闭
	if err := driver.CloseSerial(serialport); err != nil {
		err := fmt.Errorf("close serial error {%v}", err)
		log.Error(err)
		return "", err
	}

	cmdstr := fmt.Sprintf(burncmd, runtime.GOARCH, runtime.GOARCH, devport, burnfile)
	log.Infof("ArduinoUno burn cmd {%v}", cmdstr)

	cmd := exec.Command("sh", "-c", cmdstr)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Infof("ArduinoUno burn result {%v}", string(stdoutStderr))
	if err != nil {
		log.Errorf("ArduinoUno burn error {%v}", err)
		return fmt.Sprintf("%s\n\n%s", cmdstr, stdoutStderr), err
	}

	// 烧写之后进行重置
	if err := driver.Reset(devport, serialport); err != nil {
		err = fmt.Errorf("device {%v} reset error {%v}", devport, err)
		log.Error(err)
		return fmt.Sprintf("%s\n\n%s", cmdstr, stdoutStderr), err
	}

	return fmt.Sprintf("%s\n\n%s", cmdstr, stdoutStderr), nil
}

// Reset 重置
func (driver *ArduinoUno) Reset(devport string, serialport *serial.Port) error {

	/// 判断串口是否打开,如果没有打开,则重新打开串口
	if err := driver.ReOpen(serialport); err != nil {
		err := fmt.Errorf("serialport {%v} reopen error {%v}", devport, err)
		log.Error(err)
		return err
	}

	// 重置
	if err := serialport.SetDTR(false); err != nil {
		err := fmt.Errorf("ArduinoUno set dtr error {%v}", err)
		log.Error(err)
		return err
	}
	time.Sleep(time.Millisecond * 100)
	if err := serialport.SetDTR(true); err != nil {
		err := fmt.Errorf("ArduinoUno set dtr error {%v}", err)
		log.Error(err)
		return err
	}

	// 清空读取buffer
	if err := serialport.ResetInputBuffer(); err != nil {
		err := fmt.Errorf("ArduinoUno reset input buffer error {%v}", err)
		log.Error(err)
		return err
	}

	// 关闭
	if err := serialport.Close(); err != nil {
		err := fmt.Errorf("ArduinoUno close port error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}

// GetBaudRate 获取波特率
func (driver *ArduinoUno) GetBaudRate() int {
	return driver.baudrate
}

// SetBaudRate 设置波特率
func (driver *ArduinoUno) SetBaudRate(baudrate int) {
	driver.baudrate = baudrate
}

// OpenSerial 打开串口
func (driver *ArduinoUno) OpenSerial(devport string) (*serial.Port, error) {
	serialport, err := serial.Open(devport, serial.WithBaudrate(driver.baudrate), serial.WithParity(serial.NoParity), serial.WithDataBits(8), serial.WithStopBits(serial.OneStopBit))
	if err != nil {
		err = fmt.Errorf("devport {%v} open error {%v}", devport, err)
		log.Error(err)
		return nil, err
	}
	return serialport, nil
}

// CloseSerial 关闭串口
func (driver *ArduinoUno) CloseSerial(serialport *serial.Port) error {
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
func (driver *ArduinoUno) IsOpen(serialport *serial.Port) bool {
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
func (driver *ArduinoUno) ReOpen(serialport *serial.Port) error {
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
func (driver *ArduinoUno) SetResetCmd(resetcmd string) {
	return
}

// BurnEmptyProgram 烧写空程序
func (driver *ArduinoUno) BurnEmptyProgram(devport string, burncmd string, serialport *serial.Port) error {
	if output, err := driver.Burn(devport, driver.emprtprogram, burncmd, serialport); err != nil {
		err = fmt.Errorf("burn error {%v} {%v}", output, err)
		log.Error(err)
		return err
	}
	return nil
}

// SetEmptyProgram 设置空程序
func (driver *ArduinoUno) SetEmptyProgram(emptyprogram string) {
	driver.emprtprogram = emptyprogram
}

// SetWifi 设置WiFi
func (driver *ArduinoUno) SetWifi(ssid string, password string) {
}
