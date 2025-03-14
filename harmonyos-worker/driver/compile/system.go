package compile

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/database/table"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// 准备编译 system
func (cd *Driver) prepareCompileSystem(compiletable *table.CompileTable) error {

	// patch文件路径
	patchPath := filepath.Join(cd.info.Commands[compiletable.CompileType].Rootdir, "tmp.patch")

	// 写入patch文件
	if err := ioutil.WriteFile(patchPath, compiletable.FileData.Data, 0644); err != nil {
		log.Errorf("ioutil.WriteFile {%v} error {%v}", patchPath, err)
		return err
	}

	// 应用patch文件
	cmdstr := "git apply tmp.patch"
	cmd := exec.Command("sh", "-c", cmdstr)
	cmd.Dir = cd.info.Commands[compiletable.CompileType].Rootdir
	stdoutStderr, err := cmd.CombinedOutput()
	log.Debugf("%v:\n %v\n", cmdstr, string(stdoutStderr))
	if err != nil {
		err = fmt.Errorf("cmd {%v} error {%v}", cmdstr, err)
		log.Error(err)
		return err
	}

	return nil
}
