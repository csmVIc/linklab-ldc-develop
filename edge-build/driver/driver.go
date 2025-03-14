package driver

import (
	"fmt"
	"os"
	log "github.com/sirupsen/logrus"
)

// Driver 负责镜像打包
type Driver struct {
	info *DInfo
}

// DDriver 全局操作实例
var (
	DDriver *Driver
)

// New 创建实例
func New() (*Driver, error) {
	dd := Driver{info: &DInfo{}}
	// 环境变量
	if len(os.Getenv("BUILD_DOWNLOAD_URL")) < 1 {
		err := fmt.Errorf("os.Getenv(\"BUILD_DOWNLOAD_URL\") < 1 error")
		log.Error(err)
		return nil, err
	}
	dd.info.API.BuildDownloadURL = os.Getenv("BUILD_DOWNLOAD_URL")

	if len(os.Getenv("FILE_HASH")) < 1 {
		err := fmt.Errorf("os.Getenv(\"FILE_HASH\") < 1 error")
		log.Error(err)
		return nil, err
	}
	dd.info.API.FileHash = os.Getenv("FILE_HASH")

	if len(os.Getenv("TOKEN")) < 1 {
		err := fmt.Errorf("os.Getenv(\"TOKEN\") < 1 error")
		log.Error(err)
		return nil, err
	}
	dd.info.API.Token = os.Getenv("TOKEN")

	if len(os.Getenv("IMAGE_NAME")) < 1 {
		err := fmt.Errorf("os.Getenv(\"IMAGE_NAME\") < 1 error")
		log.Error(err)
		return nil, err
	}
	dd.info.Build.ImageName = os.Getenv("IMAGE_NAME")

	return &dd, nil
}
