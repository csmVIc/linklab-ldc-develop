package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) deviceupdate(client mqtt.Client, mqttmsg mqtt.Message) {

	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.Device.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var p msg.DeviceStatus
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

	// 删除失效设备
	for _, deviceid := range p.Sub {
		_, err = rdb.HDel(context.TODO(), fmt.Sprintf("devices:active:%s", username), deviceid).Result()
		if err != nil {
			log.Errorf("redis hdel error {%v}", err)
			return
		}
	}

	// 更新心跳设备
	for _, deviceid := range p.HeartBeat {
		_, err := rdb.HGet(context.TODO(), fmt.Sprintf("devices:active:%s", username), deviceid).Result()
		if err != nil {
			_, err = rdb.HSet(context.TODO(), fmt.Sprintf("devices:active:%s", username), deviceid, "").Result()
			if err != nil {
				log.Errorf("redis hset error {%v}", err)
				return
			}
		}
	}

	// 设置超时时间
	if _, err := rdb.Expire(context.TODO(), fmt.Sprintf("devices:active:%s", username), time.Duration(td.info.Device.TTL)*time.Second).Result(); err != nil {
		log.Errorf("redis expire error {%v}", err)
		return
	}
}
