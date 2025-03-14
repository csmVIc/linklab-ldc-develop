package calltest

import (
	"fmt"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

func (cd *Driver) generateSource(workdir string, index int) (string, error) {

	// 解压源代码
	sourcedir := filepath.Join(workdir, "source")
	if err := archiver.Unarchive(cd.info.Compile.SourcePath, sourcedir); err != nil {
		err := fmt.Errorf("zip unarchive {%v} to {%v} error {%v}", cd.info.Compile.SourcePath, sourcedir, err)
		log.Error(err)
		return "", err
	}

	// 修改源代码
	sourcecode := filepath.Join(sourcedir, cd.info.Compile.RandomFileName)
	if err := cd.addspace(sourcecode, index+1); err != nil {
		err := fmt.Errorf("cd.addspace {%v} {%v} error {%v}", sourcecode, index, err)
		log.Error(err)
		return "", err
	}

	// 读取文件名
	files, err := filepath.Glob(filepath.Join(sourcedir, "*"))
	if err != nil {
		err := fmt.Errorf("filepath.Glob {%v} error {%v}", filepath.Join(sourcedir, "*"), err)
		log.Error(err)
		return "", err
	}

	// 压缩文件
	outzip := filepath.Join(workdir, "outzip.zip")
	if err := archiver.Archive(files, outzip); err != nil {
		err := fmt.Errorf("zip archive {%v} to {%v} error {%v}", files, outzip, err)
		log.Error(err)
		return "", err
	}

	return outzip, nil
}
