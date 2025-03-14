package compile

import (
	"fmt"
	"linklab/device-control-v2/base-library/database/table"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// 切换branch
func (cd *Driver) switchBranch(compiletable *table.CompileTable) error {

	// 删除root目录下的所有非隐藏文件
	cmdstr := "rm -rf *"
	cmd := exec.Command("sh", "-c", cmdstr)
	cmd.Dir = cd.info.Commands[compiletable.CompileType].Rootdir
	stdoutStderr, err := cmd.CombinedOutput()
	log.Debugf("%v:\n %v\n", cmdstr, string(stdoutStderr))
	if err != nil {
		err = fmt.Errorf("cmd {%v} error {%v}", cmdstr, err)
		log.Error(err)
		return err
	}

	// 切换branch
	cmdstr = "git checkout %v"
	if compiletable.Type == "system" {
		cmdstr = fmt.Sprintf(cmdstr, compiletable.Branch)
	} else {
		cmdstr = fmt.Sprintf(cmdstr, cd.info.Commands[compiletable.CompileType].Branch)
	}
	cmd = exec.Command("sh", "-c", cmdstr)
	cmd.Dir = cd.info.Commands[compiletable.CompileType].Rootdir
	stdoutStderr, err = cmd.CombinedOutput()
	log.Debugf("%v:\n %v\n", cmdstr, string(stdoutStderr))
	if err != nil {
		err = fmt.Errorf("cmd {%v} error {%v}", cmdstr, err)
		log.Error(err)
		return err
	}

	// 恢复所有文件
	cmdstr = "git checkout ."
	cmd = exec.Command("sh", "-c", cmdstr)
	cmd.Dir = cd.info.Commands[compiletable.CompileType].Rootdir
	stdoutStderr, err = cmd.CombinedOutput()
	log.Debugf("%v:\n %v\n", cmdstr, string(stdoutStderr))
	if err != nil {
		err = fmt.Errorf("cmd {%v} error {%v}", cmdstr, err)
		log.Error(err)
		return err
	}

	// 更新到最新状态
	cmdstr = "git pull"
	cmd = exec.Command("sh", "-c", cmdstr)
	cmd.Dir = cd.info.Commands[compiletable.CompileType].Rootdir
	stdoutStderr, err = cmd.CombinedOutput()
	log.Debugf("%v:\n %v\n", cmdstr, string(stdoutStderr))
	if err != nil {
		err = fmt.Errorf("cmd {%v} error {%v}", cmdstr, err)
		log.Error(err)
		return err
	}

	return nil
}
