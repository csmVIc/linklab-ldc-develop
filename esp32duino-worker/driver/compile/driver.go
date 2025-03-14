package compile

import (
	"errors"

	"linklab/device-control-v2/base-library/parameter/msg"
	"path/filepath"
	"time"

	"github.com/mholt/archiver/v3"
	log "github.com/sirupsen/logrus"
)

// Driver 编译执行消费者
type Driver struct {
	info     *CInfo
	WaitChan chan *msg.CompileTask
}

// Driver 编译执行消费者全局实例
var (
	Cd *Driver
)

//Monitor 相应编译类型的监控协程
func (cd *Driver) Monitor() {
	for {
		parameter := <-cd.WaitChan
		cd.worker(parameter)
	}
}

// GetChannelTimeOut 获取channel超时等待时间
func (cd *Driver) GetChannelTimeOut() time.Duration {
	return time.Duration(cd.info.Channel.TimeOut) * time.Second
}

// GetSupportCompileType 获取支持的编译类型
func (cd *Driver) GetSupportCompileType() []string {
	result := make([]string, 0, len(cd.info.Commands))
	for compiletype := range cd.info.Commands {
		result = append(result, compiletype)
	}
	return result
}

// CheckCompileTypeSupportSystem 检查编译类型是否支持系统编译
func (cd *Driver) CheckCompileTypeSupportSystem(compileType string) bool {
	if cmd, isOk := cd.info.Commands[compileType]; isOk == false {
		return false
	} else {
		return cmd.SupSys
	}
}

func isEmptyDir(dir string) (bool, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil || files == nil {
		return true, err
	}
	if len(files) == 0 {
		return true, nil
	}
	return false, nil
}

// New 创建执行实例
func New(i *CInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("compile.Driver.New init info nil error")
	}
	cd := &Driver{info: i, WaitChan: make(chan *msg.CompileTask, i.Channel.Size)}

	// 工作目录初始化
	if isEmpty, err := isEmptyDir(i.Directory.Workspace); err != nil {
		log.Errorf("compile.Driver.New isEmptyDir {%v} error {%v}", i.Directory.Workspace, err)
		return nil, err
	} else if isEmpty == true {
		if err := archiver.Unarchive(i.Directory.InitZip, i.Directory.Workspace); err != nil {
			log.Errorf("compile.Driver.New unarchive {%v} to {%v} error {%v}", i.Directory.InitZip, i.Directory.Workspace, err)
			return nil, err
		}
	}

	return cd, nil
}
