package sub

import (
	"fmt"
	"linklab/device-control-v2/base-library/messenger"
	"time"

	log "github.com/sirupsen/logrus"
)

func (sd *Driver) subscriber() error {
	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		err = fmt.Errorf("get nats conn err {%v}", err)
		log.Error(err)
		return err
	}

	// 常规设备分配
	sd.callbackErr = nil
	subhandler, err := natsconn.QueueSubscribe(sd.info.Topic, "decision-maker", sd.callback)
	if err != nil {
		err = fmt.Errorf("natsconn.QueueSubscribe error {%v}", err)
		log.Error(err)
		return err
	}
	defer func() {
		err := subhandler.Unsubscribe()
		if err != nil {
			log.Errorf("nats subhandler.Unsubscribe error {%v}", err)
		}
	}()

	// 设备组分配
	gsubhandler, err := natsconn.QueueSubscribe(sd.info.GroupTopic, "decision-maker", sd.groupCallback)
	if err != nil {
		err = fmt.Errorf("natsconn.QueueSubscribe error {%v}", err)
		log.Error(err)
		return err
	}
	defer func() {
		err := gsubhandler.Unsubscribe()
		if err != nil {
			log.Errorf("nats gsubhandler.Unsubscribe error {%v}", err)
		}
	}()

	for natsconn.IsConnected() && (sd.callbackErr == nil) {
		time.Sleep(time.Second)
	}

	err = fmt.Errorf("natsconn disconnect {%v} error || nats queue subscribe {%v} error", natsconn.IsConnected(), sd.callbackErr)
	log.Error(err)
	return err
}
