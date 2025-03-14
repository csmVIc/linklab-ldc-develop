package pub

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// PubUserMsg 发布用户消息
func (pd *Driver) PubUserMsg(userMsgMap *map[string][]msg.UserMsg) error {
	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		err = fmt.Errorf("messenger.Mdriver.GetNatsConn err {%v}", err)
		log.Error(err)
		return err
	}
	for userID, userMsgList := range *userMsgMap {
		for _, userMsg := range userMsgList {
			userMsgByte, err := json.Marshal(userMsg)
			if err != nil {
				err = fmt.Errorf("json.Marshal err {%v}", err)
				log.Error(err)
				return err
			}
			err = natsconn.Publish(fmt.Sprintf(pd.info.User.Topic, userID), userMsgByte)
			if err != nil {
				err = fmt.Errorf("natsconn.Publish err {%v}", err)
				log.Error(err)
				return err
			}
		}
	}
	return nil
}
