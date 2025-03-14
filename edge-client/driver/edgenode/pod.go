package edgenode

import (
	"context"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPod 获取已有的Pod
func (ed *Driver) GetPod(namespace string, pod string) *msg.PodStatus {
	if value, isOk := ed.pods.podsMap.Load(fmt.Sprintf("%v:%v", namespace, pod)); isOk == false {
		return nil
	} else {
		return value.(*msg.PodStatus)
	}
}

// GetPods 获取已有的Pod列表
func (ed *Driver) GetPods() *msg.PodStatusList {

	res := msg.PodStatusList{
		Delete:    []msg.PodKey{},
		HeartBeat: []msg.PodStatus{},
	}

	ed.pods.podsMap.Range(func(key, value interface{}) bool {
		res.HeartBeat = append(res.HeartBeat, *value.(*msg.PodStatus))
		return true
	})

	if len(res.HeartBeat) < 1 {
		return nil
	}

	return &res
}

// getPods 获得已有的Pod列表
func (ed *Driver) getPods() *map[string]*msg.PodStatus {
	res := make(map[string]*msg.PodStatus)
	ed.pods.podsMap.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(*msg.PodStatus)
		return true
	})
	return &res
}

// lsPods 查询最新的Pod列表
func (ed *Driver) lsPods() (*map[string]*msg.PodStatus, error) {

	podlist, err := ed.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("pods list error {%v}", err)
		log.Error(err)
		return nil, err
	}

	res := make(map[string]*msg.PodStatus)
	for _, pod := range podlist.Items {
		if _, isSystem := ed.info.K8SPod.SystemNamespaces[pod.Namespace]; isSystem == true {
			// 系统Pod直接跳过
			continue
		}

		// Pod状态
		podReady := apiv1.PodCondition{
			Status: apiv1.ConditionFalse,
		}
		for _, condition := range pod.Status.Conditions {
			if condition.Type == apiv1.PodReady {
				podReady = condition
			}
		}

		elem := &msg.PodStatus{
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			NodeName:   pod.Spec.NodeName,
			Ready:      string(podReady.Status),
			CreateTime: pod.CreationTimestamp.UnixNano(),
			Containers: []msg.ContainerStatus{},
		}
		for _, cs := range pod.Status.ContainerStatuses {
			elem.Containers = append(elem.Containers, msg.ContainerStatus{
				Name:         cs.Name,
				Ready:        cs.Ready,
				RestartCount: int(cs.RestartCount),
				Image:        cs.Image,
			})
		}
		res[elem.Hash()] = elem
	}

	return &res, err
}
