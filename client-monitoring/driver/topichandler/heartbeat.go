package topichandler

import (
	"context"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/mqttclient"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) clientheartbeat(client mqtt.Client, mqttmsg mqtt.Message) {

	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	token, err := mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.HeartBeat.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	// 查找缓存
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		log.Errorf("redis get rdb error {%s}", err)
		return
	}

	// 延长token时间
	_, err = rdb.Expire(context.TODO(), token, time.Second*time.Duration(td.info.HeartBeat.TTL)).Result()
	if err != nil {
		log.Errorf("{%s} redis expire error {%v}", username, err)
		return
	}
}
