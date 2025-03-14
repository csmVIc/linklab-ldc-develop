package compile

import (
	"io/ioutil"
	"linklab/device-control-v2/base-library/database/table"
	"os"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

// 准备编译 example
func (cd *Driver) prepareCompileExample(compiletable *table.CompileTable) error {

	// 删除原先zip文件
	zipPath := cd.info.Directory.Tmp + "/" + compiletable.CompileType + ".zip"
	if _, err := os.Stat(zipPath); os.IsExist(err) {
		err = os.Remove(zipPath)
		if err != nil {
			log.Errorf("remove {%v} error {%v}", zipPath, err)
			return err
		}
	}

	// 写入zip文件
	if err := ioutil.WriteFile(zipPath, compiletable.FileData.Data, 0644); err != nil {
		log.Errorf("writefile {%v} error {%v}", zipPath, err)
		return err
	}

	// 删除原先文件
	indir := cd.info.Commands[compiletable.CompileType].Indir
	if err := os.RemoveAll(indir); err != nil {
		log.Errorf("removeall {%v} error {%v}", indir, err)
		return err
	}

	// 解压zip文件
	if err := archiver.Unarchive(zipPath, indir); err != nil {
		log.Errorf("unarchive {%v} to {%v} error {%v}", zipPath, indir, err)
		return err
	}

	return nil
}
