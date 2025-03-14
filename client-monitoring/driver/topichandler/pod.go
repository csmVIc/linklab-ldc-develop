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

func (td *Driver) podupdate(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.Pod.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var p msg.PodStatusList
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

	// 删除失效Pod
	for _, podkey := range p.Delete {
		_, err = rdb.HDel(context.TODO(), fmt.Sprintf("pods:active:%s", username), podkey.Hash()).Result()
		if err != nil {
			log.Errorf("redis hdel error {%v}", err)
			return
		}
	}

	// 更新pod节点
	for _, podstatus := range p.HeartBeat {
		podinfo := value.PodInfo{
			Name:       podstatus.Name,
			Namespace:  podstatus.Namespace,
			NodeName:   podstatus.NodeName,
			Ready:      podstatus.Ready,
			CreateTime: podstatus.CreateTime,
			Containers: []value.ContainerInfo{},
		}
		for _, containerstatus := range podstatus.Containers {
			podinfo.Containers = append(podinfo.Containers, value.ContainerInfo{
				Name:         containerstatus.Name,
				Ready:        containerstatus.Ready,
				RestartCount: containerstatus.RestartCount,
				Image:        containerstatus.Image,
			})
		}

		podinfobytes, err := json.Marshal(podinfo)
		if err != nil {
			log.Errorf("json.Marshal error {%v}", err)
			return
		}

		if _, err = rdb.HSet(context.TODO(), fmt.Sprintf("pods:active:%s", username), podinfo.Hash(), string(podinfobytes)).Result(); err != nil {
			log.Errorf("redis hset error {%v}", err)
			return
		}
	}

	// 设置超时时间
	if _, err := rdb.Expire(context.TODO(), fmt.Sprintf("pods:active:%s", username), time.Duration(td.info.Pod.TTL)*time.Second).Result(); err != nil {
		log.Errorf("redis expire error {%v}", err)
		return
	}
}
