package mqttclient

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

func (md *Driver) subinit() error {
	if md.topicToSub == nil {
		err := errors.New("mqtt topic to sub handler nil err")
		log.Error(err)
		return err
	}

	for topic, subhandler := range *md.topicToSub {
		if token := (*md.client).Subscribe(topic, subhandler.Qos, subhandler.MsgHandler); token.WaitTimeout(time.Second) && token.Error() != nil {
			log.Errorf("mqtt sub topic {%v} error", token.Error())
			return token.Error()
		}
		log.Infof("topic {%v} begin sub", topic)
	}
	return nil
}
