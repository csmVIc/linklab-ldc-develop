package aiboxnode

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetBoardFromDevPort 从串口号解析出开发板类型
func (ad *Driver) GetBoardFromDevPort(devport string) (string, error) {

	// 解析出board
	bindex := strings.LastIndex(devport, "/")
	if bindex < 0 {
		err := fmt.Errorf("devport {%v} find / index {%v} < 0", devport, bindex)
		log.Error(err)
		return "", err
	}
	eindex := strings.LastIndex(devport, "-")
	if eindex < 0 {
		err := fmt.Errorf("devport {%v} find - index {%v} < 0", devport, bindex)
		log.Error(err)
		return "", err
	}
	board := devport[bindex+1 : eindex]
	if _, isOk := ad.info.Boards[board]; isOk == false {
		err := fmt.Errorf("board {%v} not support", board)
		log.Error(err)
		return "", err
	}

	return board, nil
}
