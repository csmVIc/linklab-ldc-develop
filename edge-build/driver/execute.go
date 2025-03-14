package driver

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// Execute 镜像打包
func (dd *Driver) Execute() {

	if err := dd.prepare(); err != nil {
		err = fmt.Errorf("prepare error {%v}", err)
		log.Error(err)
		return
	}

	// 清空镜像
	defer func() {
		dd.clean()
	}()

	beginpulltime := time.Now()
	if err := dd.pull(); err != nil {
		err = fmt.Errorf("pull error {%v}", err)
		log.Error(err)
		return
	}
	beginbuildtime := time.Now()

	// 镜像打包
	if err := dd.build(); err == nil {
		beginpushtime := time.Now()
		// 镜像推送
		dd.push()
		endpushtime := time.Now()
		// 打印日志
		log.Infof("镜像拉取时间 %.2fs", beginbuildtime.Sub(beginpulltime).Seconds())
		log.Infof("镜像构建时间 %.2fs", beginpushtime.Sub(beginbuildtime).Seconds())
		log.Infof("镜像推送时间 %.2fs", endpushtime.Sub(beginpushtime).Seconds())
	}
}
