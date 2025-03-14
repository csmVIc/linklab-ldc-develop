package compile

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// pythonExecute 进行Python编译
func (cd *Driver) pythonExecute(parameter *msg.CompileTask) ([]byte, string, error) {
	// 获取源代码
	compiletable, err := cd.getCompileTable(parameter)
	if err != nil {
		log.Errorf("cd.getCompileTable error {%v}", err)
		return nil, "", err
	}

	// 切换branch
	err = cd.switchBranch(compiletable)
	if err != nil {
		log.Errorf("cd.switchBranch error {%v}", err)
		return nil, "", err
	}

	// 工作目录准备
	if compiletable.Type == "example" {
		err = cd.prepareCompileExample(compiletable)
		if err != nil {
			log.Errorf("cd.prepareCompileExample error {%v}", err)
			return nil, "", err
		}
	} else {
		err = fmt.Errorf("unsupported compile type {%v}", compiletable.Type)
		log.Error(err)
		return nil, "", err
	}

	// 拷贝文件
	indir := cd.info.Commands[parameter.CompileType].Indir
	srcfiles, err := ioutil.ReadDir(indir)
	if err != nil {
		log.Errorf("compile.Driver.pythonexecute readdir {%v} error {%v}", indir, err)
		return nil, "", err
	}
	for _, srcfile := range srcfiles {
		srcfilepath := filepath.Join(indir, srcfile.Name())
		srcfiledata, err := ioutil.ReadFile(srcfilepath)
		if err != nil {
			log.Errorf("compile.Driver.pythonexecute srcfilepath file {%v} read file error {%v}", srcfilepath, err)
			return nil, "", err
		}

		dstfilepath := filepath.Join("/app/workspace/platform/mcu/haas1000/prebuild/data", srcfile.Name())
		err = ioutil.WriteFile(dstfilepath, srcfiledata, 0644)
		if err != nil {
			log.Errorf("compile.Driver.pythonexecute dstfilepath file {%v} write file error {%v}", dstfilepath, err)
			return nil, "", err
		}
	}

	// 执行编译命令
	cmdstr := strings.ReplaceAll(cd.info.Commands[parameter.CompileType].Cmd, "boardtype", parameter.BoardType)
	cmd := exec.Command("sh", "-c", cmdstr)
	stdoutStderr, err := cmd.CombinedOutput()
	// 输出编译日志
	log.Debugf("compile out:\n %v\n", string(stdoutStderr))
	if err != nil {
		log.Errorf("compile.Driver.execute compile cmd {%v} error {%v}", cmdstr, err)
		return nil, string(stdoutStderr), err
	}

	// 检查是否出现错误
	if len(cd.info.Commands[parameter.CompileType].ErrFlag) > 0 {
		if strings.Index(string(stdoutStderr), cd.info.Commands[parameter.CompileType].ErrFlag) >= 0 {
			err = fmt.Errorf("compile stdoutStderr find {%v}", cd.info.Commands[parameter.CompileType].ErrFlag)
			log.Error(err)
			return nil, string(stdoutStderr), err
		}
	}

	outdir := strings.ReplaceAll(cd.info.Commands[parameter.CompileType].Outdir, "boardtype", parameter.BoardType)
	rregex := strings.ReplaceAll(cd.info.Commands[parameter.CompileType].Rregex, "boardtype", parameter.BoardType)
	littlefsreadpath := filepath.Join(outdir, rregex)
	littlefsreaddata, err := ioutil.ReadFile(littlefsreadpath)
	if err != nil {
		log.Errorf("compile.Driver.pythonexecute littlefsreaddata {%v} read error {%v}", littlefsreadpath, err)
		return nil, string(stdoutStderr), err
	}

	return littlefsreaddata, string(stdoutStderr), nil
}
