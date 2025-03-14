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

func (td *Driver) podresource(client mqtt.Client, mqttmsg mqtt.Message) {
	username, clientid, err := mqttclient.GetClientInfoFromTopic(mqttmsg.Topic())
	if err != nil {
		log.Errorf("mqtt get client info from topic error {%v}", err)
		return
	}

	_, err = mqttclient.MDriver.CheckClientLoginAndReplyRefuse(username, clientid, td.info.PodResource.RefuseTopic)
	if err != nil {
		log.Errorf("mqtt client check login and reply resfuse error {%v}", err)
		return
	}

	var p msg.PodResourceList
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

	// 更新
	for _, presource := range p.Pods {
		pelem := value.PodResourceInfo{
			Name:       presource.Name,
			Namespace:  presource.Namespace,
			Containers: []value.ContainerResourceInfo{},
		}
		for _, cresource := range presource.Containers {
			celem := value.ContainerResourceInfo{
				Name:   cresource.Name,
				CpuUse: cresource.CpuUse,
				MemUse: cresource.MemUse,
			}
			pelem.Containers = append(pelem.Containers, celem)
		}

		pelembytes, err := json.Marshal(pelem)
		if err != nil {
			log.Errorf("json.Marshal error {%v}", err)
			return
		}

		if _, err = rdb.HSet(context.TODO(), fmt.Sprintf("pods:resource:%s", username), pelem.Hash(), string(pelembytes)).Result(); err != nil {
			log.Errorf("redis hset error {%v}", err)
			return
		}
	}

	// 删除
	existMap := make(map[string]bool)
	for _, pod := range p.Pods {
		existMap[fmt.Sprintf("%v:%v", pod.Namespace, pod.Name)] = true
	}
	resMap, err := rdb.HGetAll(context.TODO(), fmt.Sprintf("pods:resource:%s", username)).Result()
	if err != nil {
		log.Errorf("redis hgetall error {%v}", err)
		return
	}
	for oldKey := range resMap {
		if _, isOk := existMap[oldKey]; isOk == false {
			if _, err := rdb.HDel(context.TODO(), fmt.Sprintf("pods:resource:%s", username), oldKey).Result(); err != nil {
				log.Errorf("redis hdel error {%v}", err)
				return
			}
		}
	}

	// 超时
	if _, err := rdb.Expire(context.TODO(), fmt.Sprintf("pods:resource:%s", username), time.Duration(td.info.PodResource.TTL)*time.Second).Result(); err != nil {
		log.Errorf("redis expire error {%v}", err)
		return
	}
}
