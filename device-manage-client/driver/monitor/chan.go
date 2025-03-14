package monitor

import "linklab/device-control-v2/base-library/parameter/msg"

// GetErrChan 获取错误信号信道
func (md *Driver) GetErrChan() *chan error {
	return &md.errchan
}

// GetBurnChan 获取烧写任务信道
func (md *Driver) GetBurnChan() *chan *msg.ClientBurnMsg {
	return &md.burnchan
}

// GetCmdChan 获取命令写入信道
func (md *Driver) GetCmdChan() *chan *msg.DeviceCmd {
	return &md.cmdchan
}
