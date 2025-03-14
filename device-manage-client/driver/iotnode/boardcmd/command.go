package boardcmd

import "github.com/albenik/go-serial/v2"

// BoardCommand 设备操作接口
type BoardCommand interface {
	Burn(devport string, burnfile string, burncmd string, serialport *serial.Port) (string, error)
	Reset(devport string, serialport *serial.Port) error
	GetBaudRate() int
	SetBaudRate(baudrate int)
	OpenSerial(devport string) (*serial.Port, error)
	CloseSerial(serialport *serial.Port) error
	IsOpen(serialport *serial.Port) bool
	ReOpen(serialport *serial.Port) error
	SetResetCmd(resetcmd string)
	BurnEmptyProgram(devport string, burncmd string, serialport *serial.Port) error
	SetEmptyProgram(emptyprogram string)
	SetWifi(ssid string, password string)
}
