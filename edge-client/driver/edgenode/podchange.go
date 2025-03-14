package edgenode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
)

// GetPodsChange 获取Pod列表变化
func (ed *Driver) GetPodsChange() (*msg.PodStatusList, error) {

	ed.pods.podsLock.Lock()
	defer func() {
		ed.pods.podsLock.Unlock()
	}()

	// 查询已有的Pod列表
	tmpPods := ed.getPods()
	// 查询最新的Pod列表
	nowPods, err := ed.lsPods()
	if err != nil {
		err = fmt.Errorf("ed.lsPods error {%v}", err)
		log.Error(err)
		return nil, err
	}

	res := msg.PodStatusList{
		HeartBeat: []msg.PodStatus{},
		Delete:    []msg.PodKey{},
	}
	// 查询Pod是否丢失
	for hash := range *tmpPods {
		if _, isOk := (*nowPods)[hash]; isOk == false {
			// 确认丢失
			log.Debugf("namespace {%v} pod {%v} delete", (*tmpPods)[hash].Namespace, (*tmpPods)[hash].Name)
			ed.pods.podsMap.Delete(hash)
			res.Delete = append(res.Delete, msg.PodKey{
				Name:      (*tmpPods)[hash].Name,
				Namespace: (*tmpPods)[hash].Namespace,
			})
			continue
		}
	}
	// 查询Pod是否新增
	for hash := range *nowPods {
		if _, isOk := (*tmpPods)[hash]; isOk == false {
			// 确认新增
			log.Debugf("namespace {%v} pod {%v} add", (*nowPods)[hash].Namespace, (*nowPods)[hash].Name)
			ed.pods.podsMap.Store(hash, (*nowPods)[hash])
			res.HeartBeat = append(res.HeartBeat, *(*nowPods)[hash])
			continue
		}
		// 查询Pod状态是否发生变化
		if (*(*nowPods)[hash]).Ready != (*(*tmpPods)[hash]).Ready ||
			(*(*nowPods)[hash]).NodeName != (*(*tmpPods)[hash]).NodeName {
			// 确认变化
			value, isOk := ed.pods.podsMap.Load(hash)
			if isOk == false {
				err := fmt.Errorf("pod {%v} not in podsMap error", hash)
				log.Error(err)
				return nil, err
			}
			podstatus := value.(*msg.PodStatus)
			podstatus.Ready = (*(*nowPods)[hash]).Ready
			podstatus.NodeName = (*(*nowPods)[hash]).NodeName
			podstatus.Containers = (*(*nowPods)[hash]).Containers
			log.Debugf("pod {%v} status change to {%v}", hash, *podstatus)
			res.HeartBeat = append(res.HeartBeat, *(*nowPods)[hash])
		}
	}

	return &res, nil
}
