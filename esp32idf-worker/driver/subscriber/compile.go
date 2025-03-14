package subscriber

import (
	"linklab/device-control-v2/base-library/messenger"
	"time"

	log "github.com/sirupsen/logrus"
)

func (sd *Driver) compilehandler(compileType string, supSys bool) {
	defer func() {
		sd.exitsignal = true
	}()

	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		log.Errorf("get nats conn error {%v}", err)
		return
	}

	// example 默认支持
	topicname := sd.info.ExampleMsg.Topic + "." + compileType
	subexp, err := natsconn.QueueSubscribe(topicname, sd.info.ExampleMsg.Queue, sd.msgHandler)
	if err != nil {
		log.Errorf("QueueSubscribe error {%v}", err)
		return
	}
	defer func() {
		if err := subexp.Unsubscribe(); err != nil {
			log.Errorf("Unsubscribe messenger topic {%v} error {%v}", subexp.Subject, err)
		}
	}()
	log.Debugf("sub {%v} task begin", topicname)

	// system 单独支持
	if supSys {
		topicname = sd.info.SystemMsg.Topic + "." + compileType
		subsys, err := natsconn.QueueSubscribe(topicname, sd.info.SystemMsg.Queue, sd.msgHandler)
		if err != nil {
			log.Errorf("QueueSubscribe error {%v}", err)
			return
		}
		defer func() {
			if err := subsys.Unsubscribe(); err != nil {
				log.Errorf("Unsubscribe messenger topic {%v} error {%v}", subsys.Subject, err)
			}
		}()
		log.Debugf("sub {%v} task begin", topicname)
	}

	for {
		select {
		case <-time.After(time.Second):
			if messenger.Mdriver.GetClosed() == true {
				log.Errorf("messenger.Mdriver.GetClosed true error")
				return
			}
		}
	}
}
