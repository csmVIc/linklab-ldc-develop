package topichandler

import (
	"context"
	"encoding/json"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) edgenodesetup(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.EdgeNodeSetup.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var p msg.EdgeNodeSetup
	if err := json.Unmarshal(mqttmsg.Payload(), &p); err != nil {
		log.Errorf("mqtt client payload json Unmarshal error {%v}", err)
		return
	}

	// 获取缓存
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		log.Errorf("redis get rdb error {%s}", err)
		return
	}

	// 设置
	setup := value.EdgeNodeSetupInfo{
		PodApplyURL:   p.PodApplyURL,
		PodLogURL:     p.PodLogURL,
		PodDeleteURL:  p.PodDeleteURL,
		PodExecURL:    p.PodExecURL,
		ImageBuildURL: p.ImageBuildURL,
	}
	setupvalue, err := json.Marshal(setup)
	if err != nil {
		log.Errorf("json.Marshal error {%s}", err)
		return
	}
	_, err = rdb.HSet(context.TODO(), "edgenodes:setup", username, string(setupvalue)).Result()
	if err != nil {
		log.Errorf("rdb.HSet error {%s}", err)
		return
	}
}
