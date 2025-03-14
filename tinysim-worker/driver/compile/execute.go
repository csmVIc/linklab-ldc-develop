package compile

import (
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/msg"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

// execute 进行编译
func (cd *Driver) execute(parameter *msg.CompileTask) ([]byte, string, error) {

	// 获取源代码
	binary, err := cd.getSourceCodeBinary(parameter)
	if err != nil {
		log.Errorf("getSourceCodeBinary error {%v}", err)
		return nil, "", err
	}

	// 删除原先zip文件
	zipPath := cd.info.Directory.Tmp + "/" + parameter.CompileType + ".zip"
	if _, err := os.Stat(zipPath); os.IsExist(err) {
		err = os.Remove(zipPath)
		if err != nil {
			log.Errorf("compile.Driver.execute remove {%v} error {%v}", zipPath, err)
			return nil, "", err
		}
	}

	// 写入zip文件
	if err := ioutil.WriteFile(zipPath, binary, 0644); err != nil {
		log.Errorf("compile.Driver.execute writefile {%v} error {%v}", zipPath, err)
		return nil, "", err
	}

	// 删除原先文件
	indir := cd.info.Commands[parameter.CompileType].Indir
	if err := os.RemoveAll(indir); err != nil {
		log.Errorf("compile.Driver.execute removeall {%v} error {%v}", indir, err)
		return nil, "", err
	}

	// 解压zip文件
	if err := archiver.Unarchive(zipPath, indir); err != nil {
		log.Errorf("compile.Driver.execute unarchive {%v} to {%v} error {%v}", zipPath, indir, err)
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

	// indir中的所有文件进行读取
	outfiles, err := ioutil.ReadDir(cd.info.Commands[parameter.CompileType].Indir)
	if err != nil {
		log.Errorf("compile.Driver.execute readdir {%v} error {%v}", cd.info.Commands[parameter.CompileType].Indir, err)
		return nil, string(stdoutStderr), err
	}
	compressfiles := []string{}
	for _, outfile := range outfiles {
		compressfiles = append(compressfiles, filepath.Join(cd.info.Commands[parameter.CompileType].Indir, outfile.Name()))
	}

	// 读取文件列表打包压缩
	outputzippath := filepath.Join(cd.info.Commands[parameter.CompileType].Outdir, "tinysim.zip")
	if err := archiver.Archive(compressfiles, outputzippath); err != nil {
		log.Errorf("compile.Driver.execute archiver.Archive {%v} error {%v}", outputzippath, err)
		return nil, string(stdoutStderr), err
	}

	// 读取压缩包
	outputzip, err := ioutil.ReadFile(outputzippath)
	if err != nil {
		log.Errorf("compile.Driver.execute ioutil readfile {%v} error {%v}", outputzippath, err)
		return nil, string(stdoutStderr), err
	}

	return outputzip, string(stdoutStderr), nil
}
