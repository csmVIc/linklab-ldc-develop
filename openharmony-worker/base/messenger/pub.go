package messenger

import (
	"fmt"

	"encoding/json"
)

// PubMsg 发布数据包
func (md *Driver) PubMsg(topic string, packet interface{}) error {

	conn, err := md.GetNatsConn()
	if err != nil {
		return fmt.Errorf("md.GetNatsConn error {%v}", err)
	}

	pbyte, err := json.Marshal(packet)
	if err != nil {
		return fmt.Errorf("json.Marshal error {%v}", err)
	}

	err = conn.Publish(topic, pbyte)
	if err != nil {
		return fmt.Errorf("conn.Publish {%v} error {%v}", topic, err)
	}

	return nil
}
