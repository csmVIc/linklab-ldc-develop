package compile

import (
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

// 设置编译失败状态
func (cd *Driver) setCompileError(task *msg.CompileTask, errmsg string) error {

	filter := table.CompileTableFilter{
		CompileType: task.CompileType,
		FileHash:    task.FileHash,
		BoardType:   task.BoardType,
	}
	elem := bson.M{
		"message": errmsg,
		"status":  "builderr",
	}
	updateresult, err := database.Mdriver.UpdateElem("compile", filter, elem)
	if err != nil {
		err = fmt.Errorf("database.Mdriver.UpdateElem error {%v}", err)
		log.Error(err)
		return err
	}

	if updateresult.ModifiedCount != 1 {
		err = fmt.Errorf("database.Mdriver.UpdateElem modifiedCount{%v} != 1 error", updateresult.ModifiedCount)
		log.Error(err)
		return err
	}
	return nil
}

// 设置编译成功状态
func (cd *Driver) setCompileSuccess(task *msg.CompileTask, outbin []byte) error {

	filter := table.CompileTableFilter{
		CompileType: task.CompileType,
		FileHash:    task.FileHash,
		BoardType:   task.BoardType,
	}
	elem := bson.M{
		"output": outbin,
		"status": "output",
	}
	updateresult, err := database.Mdriver.UpdateElem("compile", filter, elem)
	if err != nil {
		err = fmt.Errorf("database.Mdriver.UpdateElem error {%v}", err)
		log.Error(err)
		return err
	}

	if updateresult.ModifiedCount != 1 {
		err = fmt.Errorf("database.Mdriver.UpdateElem modifiedCount{%v} != 1 error", updateresult.ModifiedCount)
		log.Error(err)
		return err
	}

	return nil
}
