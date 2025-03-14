package iotnode

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/device-manage-client/driver/iotnode/boardcmd"
	"sync"
)

// Driver 负责设备操作
type Driver struct {
	info        *IInfo
	devices     Devices
	boardCmdMap map[string]boardcmd.BoardCommand
}

// Driver 设备操作
var (
	IDriver *Driver
)

func (id *Driver) init() error {

	// ArduinoMega2560 初始化
	if _, isOk := id.info.Commands["ArduinoMega2560"]; isOk == false {
		return errors.New("init commands not contain ArduinoMega2560")
	}
	id.boardCmdMap["ArduinoMega2560"] = &boardcmd.ArduinoMega2560{}
	id.boardCmdMap["ArduinoMega2560"].SetBaudRate(id.info.Commands["ArduinoMega2560"].BaudRate)
	id.boardCmdMap["ArduinoMega2560"].SetEmptyProgram(id.info.Commands["ArduinoMega2560"].EmptyProgram)

	// ESP32DevKitC 初始化
	if _, isOk := id.info.Commands["ESP32DevKitC"]; isOk == false {
		return errors.New("init commands not contain ESP32DevKitC")
	}
	id.boardCmdMap["ESP32DevKitC"] = &boardcmd.ESP32DevKitC{}
	id.boardCmdMap["ESP32DevKitC"].SetBaudRate(id.info.Commands["ESP32DevKitC"].BaudRate)
	id.boardCmdMap["ESP32DevKitC"].SetEmptyProgram(id.info.Commands["ESP32DevKitC"].EmptyProgram)

	// ArduinoUno 初始化
	if _, isOk := id.info.Commands["ArduinoUno"]; isOk == false {
		return errors.New("init commands not contain ArduinoUno")
	}
	id.boardCmdMap["ArduinoUno"] = &boardcmd.ArduinoUno{}
	id.boardCmdMap["ArduinoUno"].SetBaudRate(id.info.Commands["ArduinoUno"].BaudRate)
	id.boardCmdMap["ArduinoUno"].SetEmptyProgram(id.info.Commands["ArduinoUno"].EmptyProgram)

	// DeveloperKit 初始化
	if _, isOk := id.info.Commands["DeveloperKit"]; isOk == false {
		return errors.New("init commands not contain DeveloperKit")
	}
	id.boardCmdMap["DeveloperKit"] = &boardcmd.DeveloperKit{}
	id.boardCmdMap["DeveloperKit"].SetBaudRate(id.info.Commands["DeveloperKit"].BaudRate)
	id.boardCmdMap["DeveloperKit"].SetResetCmd(id.info.Commands["DeveloperKit"].Reset)
	id.boardCmdMap["DeveloperKit"].SetEmptyProgram(id.info.Commands["DeveloperKit"].EmptyProgram)

	// TelosB 初始化
	if _, isOk := id.info.Commands["TelosB"]; isOk == false {
		return errors.New("init commands not contain TelosB")
	}
	id.boardCmdMap["TelosB"] = &boardcmd.TelosB{}
	id.boardCmdMap["TelosB"].SetBaudRate(id.info.Commands["TelosB"].BaudRate)
	id.boardCmdMap["TelosB"].SetResetCmd(id.info.Commands["TelosB"].Reset)
	id.boardCmdMap["TelosB"].SetEmptyProgram(id.info.Commands["TelosB"].EmptyProgram)

	// Haas100 初始化
	if _, isOk := id.info.Commands["Haas100"]; isOk == false {
		return errors.New("init commands not contain Haas100")
	}
	id.boardCmdMap["Haas100"] = &boardcmd.Haas100{}
	id.boardCmdMap["Haas100"].SetBaudRate(id.info.Commands["Haas100"].BaudRate)
	id.boardCmdMap["Haas100"].SetEmptyProgram(id.info.Commands["Haas100"].EmptyProgram)

	// Haas100Python 初始化
	if _, isOk := id.info.Commands["Haas100Python"]; isOk == false {
		return errors.New("init commands not contain Haas100Python")
	}
	id.boardCmdMap["Haas100Python"] = &boardcmd.Haas100Python{}
	id.boardCmdMap["Haas100Python"].SetBaudRate(id.info.Commands["Haas100Python"].BaudRate)
	id.boardCmdMap["Haas100Python"].SetEmptyProgram(id.info.Commands["Haas100Python"].EmptyProgram)
	id.boardCmdMap["Haas100Python"].SetWifi(id.info.Commands["Haas100Python"].WifiSSID, id.info.Commands["Haas100Python"].WifiPassword)

	// ESP32DevKitCArduino 初始化
	if _, isOk := id.info.Commands["ESP32DevKitCArduino"]; isOk == false {
		return errors.New("init commands not contain ESP32DevKitCArduino")
	}
	id.boardCmdMap["ESP32DevKitCArduino"] = &boardcmd.ESP32DevKitCArduino{}
	id.boardCmdMap["ESP32DevKitCArduino"].SetBaudRate(id.info.Commands["ESP32DevKitCArduino"].BaudRate)
	id.boardCmdMap["ESP32DevKitCArduino"].SetEmptyProgram(id.info.Commands["ESP32DevKitCArduino"].EmptyProgram)

	// STM32F103C8 初始化
	if _, isOk := id.info.Commands["STM32F103C8"]; isOk == false {
		return errors.New("init commands not contain STM32F103C8")
	}
	id.boardCmdMap["STM32F103C8"] = &boardcmd.STM32F103C8{}
	id.boardCmdMap["STM32F103C8"].SetBaudRate(id.info.Commands["STM32F103C8"].BaudRate)
	id.boardCmdMap["STM32F103C8"].SetEmptyProgram(id.info.Commands["STM32F103C8"].EmptyProgram)

	// STM32F103C8withtest 初始化
	if _, isOk := id.info.Commands["STM32F103C8withtest"]; isOk == false {
		return errors.New("init commands not contain STM32F103C8withtest")
	}
	id.boardCmdMap["STM32F103C8withtest"] = &boardcmd.STM32F103C8withtest{}
	id.boardCmdMap["STM32F103C8withtest"].SetBaudRate(id.info.Commands["STM32F103C8withtest"].BaudRate)
	id.boardCmdMap["STM32F103C8withtest"].SetEmptyProgram(id.info.Commands["STM32F103C8withtest"].EmptyProgram)

	// STM32F103C8with4G 初始化
	if _, isOk := id.info.Commands["STM32F103C8with4G"]; isOk == false {
		return errors.New("init commands not contain STM32F103C8with4G")
	}
	id.boardCmdMap["STM32F103C8with4G"] = &boardcmd.STM32F103C8with4G{}
	id.boardCmdMap["STM32F103C8with4G"].SetBaudRate(id.info.Commands["STM32F103C8with4G"].BaudRate)
	id.boardCmdMap["STM32F103C8with4G"].SetEmptyProgram(id.info.Commands["STM32F103C8with4G"].EmptyProgram)

	// ArduinoMega2560WithHC06 初始化
	if _, isOk := id.info.Commands["ArduinoMega2560WithHC06"]; isOk == false {
		return errors.New("init commands not contain ArduinoMega2560WithHC06")
	}
	id.boardCmdMap["ArduinoMega2560WithHC06"] = &boardcmd.ArduinoMega2560WithHC06{}
	id.boardCmdMap["ArduinoMega2560WithHC06"].SetBaudRate(id.info.Commands["ArduinoMega2560WithHC06"].BaudRate)
	id.boardCmdMap["ArduinoMega2560WithHC06"].SetEmptyProgram(id.info.Commands["ArduinoMega2560WithHC06"].EmptyProgram)

	//Hi3861 初始化
	if _, isOk := id.info.Commands["Hi3861"]; isOk == false {
		return errors.New("init commands not contain Hi3861")
	}
	id.boardCmdMap["Hi3861"] = &boardcmd.Hi3861{}
	id.boardCmdMap["Hi3861"].SetBaudRate(id.info.Commands["Hi3861"].BaudRate)
	id.boardCmdMap["Hi3861"].SetEmptyProgram(id.info.Commands["Hi3861"].EmptyProgram)

	//EDURISCV64 初始化
	if _, isOk := id.info.Commands["EDURISCV64"]; isOk == false {
		return errors.New("init commands not contain EDURISCV64")
	}
	id.boardCmdMap["EDURISCV64"] = &boardcmd.EDURISCV64{}
	id.boardCmdMap["EDURISCV64"].SetBaudRate(id.info.Commands["EDURISCV64"].BaudRate)
	id.boardCmdMap["EDURISCV64"].SetEmptyProgram(id.info.Commands["EDURISCV64"].EmptyProgram)

	// SmartVilla 初始化
	if _, isOk := id.info.Commands["SmartVilla"]; isOk == false {
		return errors.New("init commands not contain SmartVilla")
	}
	id.boardCmdMap["SmartVilla"] = &boardcmd.SmartVilla{}
	id.boardCmdMap["SmartVilla"].SetBaudRate(id.info.Commands["SmartVilla"].BaudRate)
	id.boardCmdMap["SmartVilla"].SetEmptyProgram(id.info.Commands["SmartVilla"].EmptyProgram)

	// HaasEDUK1 初始化
	if _, isOk := id.info.Commands["HaasEDUK1"]; isOk == false {
		return errors.New("init commands not contain HaasEDUK1")
	}
	id.boardCmdMap["HaasEDUK1"] = &boardcmd.HaasEDUK1{}
	id.boardCmdMap["HaasEDUK1"].SetBaudRate(id.info.Commands["HaasEDUK1"].BaudRate)
	id.boardCmdMap["HaasEDUK1"].SetEmptyProgram(id.info.Commands["HaasEDUK1"].EmptyProgram)

	// ESP32DevKitCidf 初始化
	if _, isOk := id.info.Commands["ESP32DevKitCidf"]; isOk == false {
		return errors.New("init commands not contain ESP32DevKitCidf")
	}
	id.boardCmdMap["ESP32DevKitCidf"] = &boardcmd.ESP32DevKitCidf{}
	id.boardCmdMap["ESP32DevKitCidf"].SetBaudRate(id.info.Commands["ESP32DevKitCidf"].BaudRate)
	id.boardCmdMap["ESP32DevKitCidf"].SetEmptyProgram(id.info.Commands["ESP32DevKitCidf"].EmptyProgram)

	return nil


}

// New 创建设备操作
func New(i *IInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}
	id := &Driver{info: i, devices: Devices{
		deviceErrorChan: make(chan *DeviceError, i.DeviceError.ChanSize),
		devicesMap:      sync.Map{},
		devicesLock:     sync.RWMutex{},
	}, boardCmdMap: make(map[string]boardcmd.BoardCommand)}
	if err := id.init(); err != nil {
		return nil, fmt.Errorf("init error {%v}", err)
	}
	return id, nil
}
