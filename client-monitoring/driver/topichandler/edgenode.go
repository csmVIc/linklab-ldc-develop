package topichandler

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/mqttclient"
	"linklab/device-control-v2/base-library/parameter/msg"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func (td *Driver) edgenodeupdate(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.EdgeNode.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var p msg.EdgeNodeStatusList
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

	// 删除失效节点
	for _, edgenodeid := range p.Delete {
		_, err = rdb.HDel(context.TODO(), fmt.Sprintf("edgenodes:active:%s", username), edgenodeid).Result()
		if err != nil {
			log.Errorf("redis hdel error {%v}", err)
			return
		}
	}

	// 更新心跳节点
	for _, edgenodestatus := range p.HeartBeat {
		edgenodeinfo := value.EdgeNodeInfo{
			Ready:        edgenodestatus.Ready,
			Architecture: edgenodestatus.Architecture,
			OSImage:      edgenodestatus.OSImage,
			OS:           edgenodestatus.OS,
			IpAddress:    edgenodestatus.IpAddress,
			Labels:       edgenodestatus.Labels,
		}
		eninfobytes, err := json.Marshal(edgenodeinfo)
		if err != nil {
			log.Errorf("json.Marshal error {%v}", err)
			return
		}

		if _, err = rdb.HSet(context.TODO(), fmt.Sprintf("edgenodes:active:%s", username), edgenodestatus.Name, string(eninfobytes)).Result(); err != nil {
			log.Errorf("redis hset error {%v}", err)
			return
		}
	}

	// 设置超时时间
	if _, err := rdb.Expire(context.TODO(), fmt.Sprintf("edgenodes:active:%s", username), time.Duration(td.info.EdgeNode.TTL)*time.Second).Result(); err != nil {
		log.Errorf("redis expire error {%v}", err)
		return
	}
}
