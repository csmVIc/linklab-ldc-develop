package pub

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// PubClientBurnMsg 发布客户端烧写命令
func (pd *Driver) PubClientBurnMsg(clientMsgMap *map[string][]msg.ClientBurnMsg) error {
	natsconn, err := messenger.Mdriver.GetNatsConn()
	if err != nil {
		err = fmt.Errorf("messenger.Mdriver.GetNatsConn err {%v}", err)
		log.Error(err)
		return err
	}
	for clientID, clientMsgList := range *clientMsgMap {
		for _, clientMsg := range clientMsgList {
			clientMsgByte, err := json.Marshal(clientMsg)
			if err != nil {
				err = fmt.Errorf("json.Marshal err {%v}", err)
				log.Error(err)
				return err
			}
			err = natsconn.Publish(fmt.Sprintf(pd.info.Client.Burn.Topic, clientID), clientMsgByte)
			if err != nil {
				err = fmt.Errorf("natsconn.Publish err {%v}", err)
				log.Error(err)
				return err
			}
		}
	}
	return nil
}
