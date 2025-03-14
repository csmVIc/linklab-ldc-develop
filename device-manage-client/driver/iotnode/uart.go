package iotnode

import (
	"fmt"

	"github.com/albenik/go-serial/v2"
	log "github.com/sirupsen/logrus"
)

func (id *Driver) openserial(devport string) (*serial.Port, error) {

	board, err := id.GetBoardFromDevPort(devport)
	if err != nil {
		err = fmt.Errorf("get board from devport error {%v}", err)
		log.Error(err)
		return nil, err
	}

	serialport, err := id.boardCmdMap[board].OpenSerial(devport)
	if err != nil {
		err = fmt.Errorf("devport {%v} open error {%v}", devport, err)
		log.Error(err)
		return nil, err
	}

	return serialport, nil
}
