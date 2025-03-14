package driver

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

func (dd *Driver) prepare() error {

	// 文件下载
	zipData, err := dd.buildDownload()
	if err != nil {
		log.Errorf("buildDownload error {%v}", err)
		return err
	}

	zipPath := filepath.Join("/app/tmp", "build.zip")
	if err := ioutil.WriteFile(zipPath, zipData, 0644); err != nil {
		log.Errorf("writefile {%v} error {%v}", zipPath, err)
		return err
	}

	// 文件解压
	unZipPath := "/app/workspace"
	if err := archiver.Unarchive(zipPath, unZipPath); err != nil {
		log.Errorf("unarchive {%v} to {%v} error {%v}", zipPath, unZipPath, err)
		return err
	}

	return nil
}
