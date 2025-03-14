package calltest

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (cd *Driver) addspace(filepath string, length int) error {

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		err := fmt.Errorf("failed opening file {%v} error {%v}", filepath, err)
		log.Error(err)
		return err
	}
	defer file.Close()

	sb := strings.Builder{}
	for i := 0; i < length; i++ {
		sb.WriteString(" ")
	}

	len, err := file.WriteString(sb.String())
	if err != nil {
		err := fmt.Errorf("failed writing to file {%v} error {%v}", len, err)
		log.Error(err)
		return err
	}
	log.Debugf("write to file {%v} len {%v}", filepath, len)
	return nil
}
