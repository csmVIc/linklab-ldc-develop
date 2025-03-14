package compile

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/database/table"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// 准备编译 system
func (cd *Driver) prepareCompileSystem(compiletable *table.CompileTable) error {

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

	// 解压zip文件
	rootdir := cd.info.Commands[compiletable.CompileType].Rootdir
	cmdstr := fmt.Sprintf("unzip -o %v -d %v", zipPath, rootdir)
	cmd := exec.Command("sh", "-c", cmdstr)
	// cmd.Dir = rootdir
	stdoutStderr, err := cmd.CombinedOutput()
	log.Debugf("%v:\n %v\n", cmdstr, string(stdoutStderr))
	if err != nil {
		log.Errorf("cmd {%v} error {%v}", cmd, err)
		return err
	}

	return nil
}
