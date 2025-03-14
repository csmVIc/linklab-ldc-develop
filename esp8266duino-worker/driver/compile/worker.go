package compile

import (
	"os"
	"time"

	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// Monitor 监控编译任务提交请求,并启动相应的worker
func (cd *Driver) worker(parameter *msg.CompileTask) error {

	log.Debugf("begin compile {%v}", *parameter)

	begincompiletime := time.Now()
	outbin, outmsg, err := cd.execute(parameter)
	endcompiletime := time.Now()

	// nil check
	errmsg := ""
	if err != nil {
		errmsg = err.Error()
		log.Errorf("cd.execute error {%v}", errmsg)
	}

	// 日志记录
	tags := map[string]string{
		"filehash":    parameter.FileHash,
		"compiletype": parameter.CompileType,
		"boardtype":   parameter.BoardType,
		"podname":     os.Getenv("POD_NAME"),
		"nodename":    os.Getenv("NODE_NAME"),
	}
	fields := map[string]interface{}{
		"begincompiletime": begincompiletime.Format(time.RFC3339),
		"endcompiletime":   endcompiletime.Format(time.RFC3339),
		"compiletime":      endcompiletime.Sub(begincompiletime).Seconds(),
		"outmsg":           outmsg,
		"outerr":           errmsg,
	}
	err = logger.Ldriver.WriteLog("compilelog", tags, fields)
	if err != nil {
		log.Errorf("logger.Ldriver.WriteLog error {%v}", err)
		return err
	}

	if len(errmsg) < 1 {
		// 编译成功
		err = cd.setCompileSuccess(parameter, outbin)
		if err != nil {
			log.Errorf("cd.setCompileSuccess error {%v}", err)
			return err
		}
	} else {
		// 编译失败
		err = cd.setCompileError(parameter, outmsg)
		if err != nil {
			log.Errorf("cd.setCompileError error {%v}", err)
			return err
		}
	}

	log.Debugf("end compile {%v}", *parameter)
	return nil
}
