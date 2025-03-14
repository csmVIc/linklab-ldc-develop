package monitor

// GetErrChan 获取错误信号信道
func (md *Driver) GetErrChan() *chan error {
	return &md.errchan
}

// GetPodApplyChan 获取Pod部署信道
// func (md *Driver) GetPodApplyChan() *chan *msg.PodApply {
// 	return &md.podapplychan
// }
