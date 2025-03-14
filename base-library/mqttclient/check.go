package mqttclient

import (
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// CheckClientLoginAndReplyRefuse 检查客户端的登录状态
func (md *Driver) CheckClientLoginAndReplyRefuse(username string, clientid string, topic string) (string, error) {

	var token string
	var err error
	if token, err = auth.CheckClientIDAndGetToken(username, clientid); err != nil {
		log.Errorf("check client token and id error {%v}", err)

		reply := msg.ReplyMsg{
			Code: -1,
			Msg:  err.Error(),
			Data: nil,
		}
		if err := md.PubMsg(fmt.Sprintf(topic, username, clientid), 0, reply); err != nil {
			log.Errorf("{%s} publish token error {%v}", username, err)
			return "", err
		}
		return "", err
	}

	return token, nil
}
