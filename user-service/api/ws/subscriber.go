package ws

import (
	"fmt"
	"linklab/device-control-v2/base-library/messenger"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func subscriber(msgchan *chan *nats.Msg, exitchan *chan bool, topic string) error {

	conn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		log.Errorf("get nats conn error {%v}", err)
		return err
	}

	sub, err := conn.Subscribe(topic, func(msg *nats.Msg) {
		*msgchan <- msg
		log.Infof("sub topic {%v} msg {%v}", topic, string(msg.Data))
	})

	if err != nil {
		log.Errorf("subscribe messenger topic {%v} error {%v}", topic, err)
		return err
	}

	defer func() {
		if err := sub.Unsubscribe(); err != nil {
			log.Errorf("unsubscribe messenger topic {%v} error {%v}", topic, err)
		}
	}()

	for true {
		select {
		case <-*exitchan:
			log.Infof("subscriber {%v} recv exit signal", topic)
			return nil
		case <-time.After(time.Second):
			if conn.IsConnected() == false {
				err := fmt.Errorf("subscriber {%v} disconnect error", topic)
				log.Error(err)
				return err
			}
		}
	}
	return nil
}
