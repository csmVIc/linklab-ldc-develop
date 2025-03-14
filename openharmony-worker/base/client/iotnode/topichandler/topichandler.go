package topichandler

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// Driver 负责mqtt消息处理
type Driver struct {
	info      *TInfo
	errchan   *chan error
	burnchan  *chan *msg.ClientBurnMsg
	tokenchan *chan string
	cmdchan   *chan *msg.DeviceCmd
	// podapplychan *chan *msg.PodApply
	topicMap   *map[string]TopicInfo
	topicToSub *map[string]mqttclient.SubTopic
}

// TDriver mqtt消息处理全局实例
var (
	TDriver *Driver
)

func (td *Driver) init() error {

	// 获取username和clientid
	username, err := mqttclient.MDriver.GetUserName()
	if err != nil {
		log.Errorf("mqttclient.MDriver.GetUserName error {%v}", err)
		return err
	}
	clientid, err := mqttclient.MDriver.GetClientID()
	if err != nil {
		log.Errorf("mqttclient.MDriver.GetClientID error {%v}", err)
		return err
	}

	// 配置信息检查更新
	tm := map[string]TopicInfo{}
	for key, elem := range td.info.Topics {
		if len(elem.Pub) > 0 {
			elem.Pub = fmt.Sprintf(elem.Pub, username, clientid)
		}
		if len(elem.Sub) > 0 {
			elem.Sub = fmt.Sprintf(elem.Sub, username, clientid)
		}
		if len(elem.Refuse) > 0 {
			elem.Refuse = fmt.Sprintf(elem.Refuse, username, clientid)
		}
		tm[key] = elem
	}
	td.topicMap = &tm

	// 创建mqtt订阅关系
	ts := make(map[string]mqttclient.SubTopic)
	ts[(*td.topicMap)["heartbeat"].Refuse] = mqttclient.SubTopic{
		MsgHandler: td.heartbeatrefuse,
		Qos:        0,
	}
	ts[(*td.topicMap)["authtoken"].Sub] = mqttclient.SubTopic{
		MsgHandler: td.authtokensub,
		Qos:        2,
	}
	ts[(*td.topicMap)["authtoken"].Refuse] = mqttclient.SubTopic{
		MsgHandler: td.authtokenrefuse,
		Qos:        0,
	}
	if len((*td.topicMap)["deviceupdate"].Refuse) > 0 {
		ts[(*td.topicMap)["deviceupdate"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.deviceupdaterefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["devicelog"].Refuse) > 0 {
		ts[(*td.topicMap)["devicelog"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.devicelogrefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["endrun"].Refuse) > 0 {
		ts[(*td.topicMap)["endrun"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.endrunrefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["deviceburn"].Sub) > 0 {
		ts[(*td.topicMap)["deviceburn"].Sub] = mqttclient.SubTopic{
			MsgHandler: td.deviceburnsub,
			Qos:        2,
		}
	}
	// 需要烧写进行该初始化
	if len((*td.topicMap)["deviceburn"].Refuse) > 0 {
		ts[(*td.topicMap)["deviceburn"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.burnresultrefuse,
			Qos:        0,
		}
	}
	// 需要执行进行该初始化
	if len((*td.topicMap)["execerr"].Refuse) > 0 {
		ts[(*td.topicMap)["execerr"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.execerrrefuse,
			Qos:        0,
		}
	}
	// 需要命令写入进行该初始化
	if len((*td.topicMap)["cmdwrite"].Sub) > 0 {
		ts[(*td.topicMap)["cmdwrite"].Sub] = mqttclient.SubTopic{
			MsgHandler: td.cmdwritesub,
			Qos:        2,
		}
	}
	if len((*td.topicMap)["edgenodeupdate"].Refuse) > 0 {
		ts[(*td.topicMap)["edgenodeupdate"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.edgenodeupdaterefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["podupdate"].Refuse) > 0 {
		ts[(*td.topicMap)["podupdate"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.podupdaterefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["edgenoderesource"].Refuse) > 0 {
		ts[(*td.topicMap)["edgenoderesource"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.edgenoderesourcerefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["podresource"].Refuse) > 0 {
		ts[(*td.topicMap)["podresource"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.podresourcerefuse,
			Qos:        0,
		}
	}
	if len((*td.topicMap)["edgenodesetup"].Refuse) > 0 {
		ts[(*td.topicMap)["edgenodesetup"].Refuse] = mqttclient.SubTopic{
			MsgHandler: td.edgenodesetuprefuse,
			Qos:        0,
		}
	}

	// if len((*td.topicMap)["podapply"].Sub) > 0 {
	// 	ts[(*td.topicMap)["podapply"].Sub] = mqttclient.SubTopic{
	// 		MsgHandler: td.podapplysub,
	// 		Qos:        2,
	// 	}
	// }
	// if len((*td.topicMap)["podapply"].Refuse) > 0 {
	// 	ts[(*td.topicMap)["podapply"].Refuse] = mqttclient.SubTopic{
	// 		MsgHandler: td.podapplyresultresfuse,
	// 		Qos:        0,
	// 	}
	// }

	td.topicToSub = &ts
	return nil
}

// New 创建mqtt消息处理实例
func New(i *TInfo) (*Driver, error) {

	if i == nil {
		err := errors.New("init info i nil error")
		log.Error(err)
		return nil, err
	}

	td := &Driver{info: i, topicMap: nil, topicToSub: nil}
	if err := td.init(); err != nil {
		err = fmt.Errorf("td.init error {%v}", err)
		log.Error(err)
		return nil, err
	}

	return td, nil
}

// GetTopicSubHandler 获取mqtt消息处理
func (td *Driver) GetTopicSubHandler() *map[string]mqttclient.SubTopic {
	return td.topicToSub
}
