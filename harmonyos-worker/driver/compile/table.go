package compile

import (
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

func (cd *Driver) getCompileTable(parameter *msg.CompileTask) (*table.CompileTable, error) {

	filter := &table.CompileTableFilter{
		CompileType: parameter.CompileType,
		BoardType:   parameter.BoardType,
		FileHash:    parameter.FileHash,
	}
	result := &table.CompileTable{}
	err := database.Mdriver.FindOneElem("compile", filter, result)
	if err != nil {
		err = fmt.Errorf("database.Mdriver.FindOneElem error {%v}", err)
		log.Error(err)
		return nil, err
	}

	if result.Status != "input" {
		err = fmt.Errorf("database.Mdriver.FindOneElem find status {%v} != {input}", result.Status)
		log.Error(err)
		return nil, err
	}

	return result, nil
}
