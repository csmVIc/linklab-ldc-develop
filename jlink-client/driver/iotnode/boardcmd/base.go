package boardcmd

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

// 获取JLink编号
func GetJLinkSerialNumber(devport string) (int, error) {

	index := strings.LastIndex(devport, "-")
	if index < 0 {
		return -1, fmt.Errorf("parse {%v} serial number error", devport)
	}

	num, err := strconv.Atoi(devport[index+1:])
	if err != nil {
		log.Error(err)
		return -1, err
	}

	return num, nil
}
