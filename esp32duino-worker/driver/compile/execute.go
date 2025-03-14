package compile

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// execute 进行编译
func (cd *Driver) execute(parameter *msg.CompileTask) ([]byte, string, error) {

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
	} else if compiletable.Type == "system" {
		err = cd.prepareCompileSystem(compiletable)
		if err != nil {
			log.Errorf("cd.prepareCompileSystem error {%v}", err)
			return nil, "", err
		}
	} else {
		err = fmt.Errorf("unsupported compile type {%v}", compiletable.Type)
		log.Error(err)
		return nil, "", err
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

	// 读取输出结果
	outdir := strings.ReplaceAll(cd.info.Commands[parameter.CompileType].Outdir, "boardtype", parameter.BoardType)
	files, err := ioutil.ReadDir(outdir)
	if err != nil {
		log.Errorf("compile.Driver.execute read out dir {%v} error {%v}", outdir, err)
		return nil, string(stdoutStderr), err
	}
	rregex := strings.ReplaceAll(cd.info.Commands[parameter.CompileType].Rregex, "boardtype", parameter.BoardType)
	outmatch, err := regexp.Compile(rregex)
	if err != nil {
		log.Errorf("compile.Driver.execute regexp compile {%v} error {%v}", rregex, err)
		return nil, string(stdoutStderr), err
	}
	patmatch, err := regexp.Compile("^.*partitions.*")
	if err != nil {
		log.Errorf("compile.Driver.execute regexp compile {%v} error {%v}", "^.*partitions.*", err)
		return nil, string(stdoutStderr), err
	}

	isOk := false
	outfname := ""
	for _, file := range files {
		isOk = outmatch.MatchString(file.Name()) && (!patmatch.MatchString(file.Name()))
		if isOk {
			log.Infof("compile.Driver.execute match out file name {%v}", file.Name())
			outfname = file.Name()
			break
		}
		log.Infof("compile.Driver.execute unmatch out file name {%v}", file.Name())
	}

	if isOk {
		outbin, err := ioutil.ReadFile(outdir + "/" + outfname)
		if err != nil {
			log.Errorf("compile.Driver.execute ioutil readfile {%v} error {%v}", outdir+"/"+outfname, err)
			return nil, string(stdoutStderr), err
		}
		return outbin, string(stdoutStderr), nil
	}

	err = fmt.Errorf("compile.Driver.execute can not find output file {%v}", outmatch.String())
	log.Error(err)
	return nil, string(stdoutStderr), err
}
