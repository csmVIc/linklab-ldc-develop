package subscriber

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/arduinomega-virtual-worker/driver/compile"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
)

func (sd *Driver) msgHandler(natmsg *nats.Msg) {

	respond := msg.ReplyMsg{
		Code: 0,
		Msg:  "success",
		Data: nil,
	}

	task := msg.CompileTask{}
	if err := json.Unmarshal(natmsg.Data, &task); err != nil {
		err := fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		respond.Code = -1
		respond.Msg = err.Error()
	} else {
		log.Infof("topic {%v} sub task {%v}", natmsg.Subject, task)
		timeout := compile.Cd.GetChannelTimeOut()
		select {
		case compile.Cd.WaitChan <- &task:
			log.Infof("task {%v} enter compiling queue success", task)
		case <-time.After(timeout):
			err := fmt.Errorf("task {%v} enter compiling queue timeout %v error", task, timeout)
			log.Error(err)
			respond.Code = -1
			respond.Msg = err.Error()
		}
	}

	// 回复消息
	respondBytes, err := json.Marshal(respond)
	if err != nil {
		log.Errorf("json.Marshal error {%v}", err)
		sd.exitsignal = true
		return
	}
	if err = natmsg.Respond(respondBytes); err != nil {
		log.Errorf("natmsg.Respond error {%v}", err)
		sd.exitsignal = true
		return
	}
}
