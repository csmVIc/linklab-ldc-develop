package sub

import (
	"fmt"
	"linklab/device-control-v2/base-library/messenger"
	"time"

	log "github.com/sirupsen/logrus"
)

func (sd *Driver) subscribe() error {

	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		err = fmt.Errorf("get nats conn err {%v}", err)
		log.Error(err)
		return err
	}

	sd.callbackErr = nil

	// pod部署
	podApplySubHandler, err := natsconn.QueueSubscribe(sd.info.Topic.PodApply, "edge-allocater", sd.podApplyCallback)
	if err != nil {
		err = fmt.Errorf("natsconn.QueueSubscribe error {%v}", err)
		log.Error(err)
		return err
	}
	defer func() {
		err := podApplySubHandler.Unsubscribe()
		if err != nil {
			log.Errorf("nats podApplySubHandler.Unsubscribe error {%v}", err)
		}
	}()

	// 镜像构建
	imageBuildSubHandler, err := natsconn.QueueSubscribe(sd.info.Topic.ImageBuild, "edge-allocater", sd.imageBuildCallback)
	if err != nil {
		err = fmt.Errorf("natsconn.QueueSubscribe error {%v}", err)
		log.Error(err)
		return err
	}
	defer func() {
		err := imageBuildSubHandler.Unsubscribe()
		if err != nil {
			log.Errorf("nats imageBuildSubHandler.Unsubscribe error {%v}", err)
		}
	}()

	for natsconn.IsConnected() && (sd.callbackErr == nil) {
		time.Sleep(time.Second)
	}

	err = fmt.Errorf("natsconn disconnect {%v} error || nats queue subscribe {%v} error", natsconn.IsConnected(), sd.callbackErr)
	log.Error(err)
	return err
}
