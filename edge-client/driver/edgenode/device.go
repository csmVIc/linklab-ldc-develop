package edgenode

import (
	"context"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetEdgeNodes 获取已有的设备列表
func (ed *Driver) GetEdgeNodes() *msg.EdgeNodeStatusList {

	res := msg.EdgeNodeStatusList{
		Delete:    []string{},
		HeartBeat: []msg.EdgeNodeStatus{},
	}

	ed.edgenodes.edgeNodesMap.Range(func(key, value interface{}) bool {
		res.HeartBeat = append(res.HeartBeat, *value.(*msg.EdgeNodeStatus))
		return true
	})

	if len(res.HeartBeat) < 1 {
		return nil
	}

	return &res
}

// getEdgeNodes 获取已有的设备列表
func (ed *Driver) getEdgeNodes() *map[string]*msg.EdgeNodeStatus {
	res := make(map[string]*msg.EdgeNodeStatus)
	ed.edgenodes.edgeNodesMap.Range(func(key, value interface{}) bool {
		res[key.(string)] = value.(*msg.EdgeNodeStatus)
		return true
	})
	return &res
}

// lsEdgeNodes 查询最新的设备列表
func (ed *Driver) lsEdgeNodes() (*map[string]*msg.EdgeNodeStatus, error) {
	nodelist, err := ed.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("nodes list error {%v}", err)
		log.Error(err)
		return nil, err
	}

	res := make(map[string]*msg.EdgeNodeStatus)
	for _, node := range nodelist.Items {
		if _, isOk := node.Labels["node-role.kubernetes.io/master"]; isOk {
			// 非边缘节点直接跳过
			continue
		}

		// 节点状态
		nodeReady := apiv1.NodeCondition{
			Status: apiv1.ConditionFalse,
		}
		for _, condition := range node.Status.Conditions {
			if condition.Type == apiv1.NodeReady {
				nodeReady = condition
			}
		}

		// IP地址
		ipAddress := "None"
		for _, address := range node.Status.Addresses {
			if address.Type == apiv1.NodeInternalIP {
				ipAddress = address.Address
				break
			}
		}

		elem := &msg.EdgeNodeStatus{
			Name:         node.Name,
			Ready:        string(nodeReady.Status),
			Architecture: node.Status.NodeInfo.Architecture,
			OSImage:      node.Status.NodeInfo.OSImage,
			OS:           node.Status.NodeInfo.OperatingSystem,
			IpAddress:    ipAddress,
			Labels:       node.Labels,
		}
		res[node.Name] = elem
	}

	return &res, nil
}
