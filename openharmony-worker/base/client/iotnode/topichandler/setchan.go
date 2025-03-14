package topichandler

import "linklab/device-control-v2/base-library/parameter/msg"

func (td *Driver) SetErrChan(ec *chan error) {
	td.errchan = ec
}

func (td *Driver) SetBurnChan(bc *chan *msg.ClientBurnMsg) {
	td.burnchan = bc
}

func (td *Driver) SetTokenChan(tc *chan string) {
	td.tokenchan = tc
}

func (td *Driver) SetCmdChan(cc *chan *msg.DeviceCmd) {
	td.cmdchan = cc
}

// func (td *Driver) SetPodApplyChan(pac *chan *msg.PodApply) {
// 	td.podapplychan = pac
// }
