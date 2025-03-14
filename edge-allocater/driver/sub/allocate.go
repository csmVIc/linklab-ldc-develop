package sub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/tool"
	"math/rand"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (sd *Driver) allocate(userid string, podname string, nodeselector map[string]string) (string, string, error) {

	// 获取用户的租户信息
	utenantid, err := cache.Cdriver.GetUserTenantID(userid)
	if err != nil {
		err := fmt.Errorf("cache.Cdriver.GetUserTenantID {%v} error {%v}", userid, err)
		log.Error(err)
		return "", "", err
	}

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return "", "", err
	}

	// 遍历所有Pods
	workJetson := []string{}
	tkeys, err := rdb.Keys(context.TODO(), "pods:active:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		return "", "", err
	}

	// 遍历所有边缘客户端
	availClients := []string{}
	edgeClients, err := rdb.Keys(context.TODO(), "edgenodes:active:*").Result()
	if err != nil {
		err = fmt.Errorf("rdb.Keys error {%s}", err)
		log.Error(err)
		return "", "", err
	}

	// 如果nodeselector为空，那么指定默认设备类型为树莓派
	if nodeselector == nil || len(nodeselector) < 1 {
		nodeselector = map[string]string{
			"linklab.edgetype": "RaspberryPi4BExtend",
		}
	}

	// 获取正在工作的jetsonnano设备
	if nodeselector["linklab.edgetype"] == "jetsonnano"{
		for _, tkey := range tkeys {
			podmap, err := rdb.HGetAll(context.TODO(), tkey).Result()
			if err != nil {
				err = fmt.Errorf("redis hgetall error {%v}", err)
				log.Error(err)
				return "", "", err
			}
	
			for _, podvaluestr := range podmap {
	
				podvalue := &value.PodInfo{}
				if err := json.Unmarshal([]byte(podvaluestr), podvalue); err != nil {
					err = fmt.Errorf("json.Unmarshal error {%v}", err)
					log.Error(err)
					return "", "", err
				}

				if  strings.Split(podvalue.NodeName, "-")[0] == "jetsonnano"{
					workJetson = append(workJetson, podvalue.NodeName)
				}
			}
		}
	}
	
	for _, edgeClient := range edgeClients {
		clientid := edgeClient[strings.LastIndex(edgeClient, ":")+1:]
		// 查询客户端所属的租户ID
		clientTenantID, err := cache.Cdriver.GetClientTenantID(clientid)
		if err != nil {
			err = fmt.Errorf("cache.Cdriver.GetClientTenantID {%v} error {%s}", clientid, err)
			log.Error(err)
			return "", "", err
		}

		// 检查租户信息是否符合
		if _, isOk := clientTenantID[utenantid]; !isOk {
			continue
		}

		// 检查是否拥有所需类型的节点
		if len(nodeselector) > 0 {
			edgenodes, err := rdb.HGetAll(context.TODO(), edgeClient).Result()
			if err != nil {
				err = fmt.Errorf("redis hgetall error {%v}", err)
				log.Error(err)
				return "", "", err
			}

			isOk := false
			for _, envaluestr := range edgenodes {
				envalue := &value.EdgeNodeInfo{}
				if err := json.Unmarshal([]byte(envaluestr), envalue); err != nil {
					err = fmt.Errorf("json.Unmarshal error {%v}", err)
					log.Error(err)
					return "", "", err
				}
				// 检查是否有raspberry和jetsonnano节点，检查jetsonnano节点是否在工作中
				if tool.MapAIncludeMapB(envalue.Labels, nodeselector) || nodeselector["linklab.edgetype"] == "jetsonnano"{
					if strings.Split(envalue.Labels["kubernetes.io/hostname"], "-")[0] == "jetsonnano"{
						currNodename := envalue.Labels["kubernetes.io/hostname"]
						isfree := true
						for _, element := range workJetson{
							if currNodename == element{
								isfree = false
								break
							}
						}
						if isfree == true && nodeselector["linklab.edgetype"] == "jetsonnano"{
							nodeselector["linklab.edgetype"] = currNodename
							isOk = true
							break
						}
					}else{
						if nodeselector["linklab.edgetype"] != "jetsonnano"{
							isOk = true
							break
						}	
					}
				}
			}
			// 未满足要求
			if isOk == false {
				continue
			}
		}
		// 满足所有约束条件
		availClients = append(availClients, clientid)
	}

	if len(availClients) < 1 {
		return "", "", errors.New("available clients length < 1 error")
	}

	// 检查是否已分配过该Pod
	namespace := tool.CreateMD5(userid)
	for _, edgeClient := range availClients {
		exist, err := rdb.HExists(context.TODO(), fmt.Sprintf("pods:active:%v", edgeClient), fmt.Sprintf("%v:%v", namespace, podname)).Result()
		if err != nil {
			err = fmt.Errorf("rdb.HExists error {%s}", err)
			log.Error(err)
			return "", "", err
		}
		// 已分配过，直接返回
		if exist {
			return edgeClient, "", nil
		}
	}

	return availClients[rand.Intn(len(availClients))], nodeselector["linklab.edgetype"],  nil
}
