package mqttclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// PubMsg 发布消息
func (md *Driver) PubMsg(topic string, qos byte, msg interface{}) error {

	if md.client == nil {
		return errors.New("mqtt client nil err")
	}

	msgByte, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json Marshal error {%v}", err)
	}

	if ack := (*md.client).Publish(topic, qos, false, msgByte); ack.WaitTimeout(time.Duration(md.info.Publish.TimeOut)*time.Second) && ack.Error() != nil {
		return fmt.Errorf("mqtt publish {%v} error {%v}", topic, ack.Error())
	}

	return nil
}

// PubMsgByte 发布消息
func (md *Driver) PubMsgByte(topic string, qos byte, msgByte []byte) error {

	if md.client == nil {
		return errors.New("mqtt client nil err")
	}

	if ack := (*md.client).Publish(topic, qos, false, msgByte); ack.WaitTimeout(time.Duration(md.info.Publish.TimeOut)*time.Second) && ack.Error() != nil {
		return fmt.Errorf("mqtt publish {%v} error {%v}", topic, ack.Error())
	}

	return nil
}
