package subscriber

import (
	"linklab/device-control-v2/base-library/messenger"
	"time"

	log "github.com/sirupsen/logrus"
)

func (sd *Driver) compilehandler(compileType string) {
	defer func() {
		sd.exitsignal = true
	}()

	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		log.Errorf("get nats conn error {%v}", err)
		return
	}

	topicname := sd.info.Topic + "." + compileType
	sub, err := natsconn.QueueSubscribe(topicname, "worker", sd.msgHandler)
	if err != nil {
		log.Errorf("QueueSubscribe error {%v}", err)
		return
	}
	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Errorf("Unsubscribe messenger topic {%v} error {%v}", sub.Subject, err)
		}
	}()

	log.Debugf("sub {%v} task begin", topicname)

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
