package messenger

import (
	"encoding/json"
	"fmt"
	"time"
)

// RequestMsg 请求并获取回复数据包
func (md *Driver) RequestMsg(topic string, packet interface{}, timeout time.Duration, reply interface{}) error {

	conn, err := md.GetNatsConn()
	if err != nil {
		return fmt.Errorf("md.GetNatsConn error {%v}", err)
	}

	pbyte, err := json.Marshal(packet)
	if err != nil {
		return fmt.Errorf("json.Marshal error {%v}", err)
	}

	replyMsg, err := conn.Request(topic, pbyte, timeout)
	if err != nil {
		return fmt.Errorf("conn.Request {%v} error {%v}", topic, err)
	}

	err = json.Unmarshal(replyMsg.Data, reply)
	if err != nil {
		return fmt.Errorf("json.Unmarshal {%v} error {%v}", string(replyMsg.Data), err)
	}

	return nil
}
