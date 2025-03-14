package edgenode

import (
	"context"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (ed *Driver) GetEdgeNodeResourceList() (*msg.EdgeNodeResourceList, error) {

	nodelist, err := ed.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("nodes list error {%v}", err)
		log.Error(err)
		return nil, err
	}

	nmap := make(map[string]*msg.EdgeNodeResource)
	for _, node := range nodelist.Items {
		if _, isOk := node.Labels["node-role.kubernetes.io/master"]; isOk {
			// 非边缘节点直接跳过
			continue
		}

		elem := &msg.EdgeNodeResource{
			Name:         node.Name,
			CpuAll:       node.Status.Allocatable.Cpu().MilliValue(),
			MemAll:       node.Status.Allocatable.Memory().MilliValue(),
			NvidiaGpuAll: 0,
		}

		// GPU
		if gpuRes, isOk := node.Status.Allocatable["nvidia.com/gpu"]; isOk {
			elem.NvidiaGpuAll = gpuRes.Value()
		}

		// 查询边缘节点的资源消耗
		edgemetrics := &v1beta1.NodeMetrics{}
		err = ed.clientset.RESTClient().Get().AbsPath(fmt.Sprintf("apis/metrics.k8s.io/v1beta1/nodes/%v", node.Name)).Do(context.TODO()).Into(edgemetrics)
		if err != nil {
			err = fmt.Errorf("get metrics error {%v}", err)
			log.Error(err)
			return nil, err
		}

		elem.CpuUse = edgemetrics.Usage.Cpu().MilliValue()
		elem.MemUse = edgemetrics.Usage.Memory().MilliValue()

		nmap[node.Name] = elem
	}

	res := msg.EdgeNodeResourceList{
		EdgeNodes: []msg.EdgeNodeResource{},
	}
	for _, elem := range nmap {
		res.EdgeNodes = append(res.EdgeNodes, *elem)
	}
	return &res, nil
}

func (ed *Driver) GetPodResourceList() (*msg.PodResourceList, error) {

	// 获取所有namespace
	nslist, err := ed.lsNamespace()
	if err != nil {
		err = fmt.Errorf("ed.lsNamespace error {%v}", err)
		log.Error(err)
		return nil, err
	}

	res := &msg.PodResourceList{
		Pods: []msg.PodResource{},
	}
	for _, namespace := range nslist {
		podmetricslist := &v1beta1.PodMetricsList{}
		if err = ed.clientset.RESTClient().Get().AbsPath(fmt.Sprintf("apis/metrics.k8s.io/v1beta1/namespaces/%v/pods", namespace)).Do(context.TODO()).Into(podmetricslist); err != nil {
			err = fmt.Errorf("get metrics error {%v}", err)
			log.Error(err)
			return nil, err
		}

		for _, podmetrics := range podmetricslist.Items {
			pelem := msg.PodResource{
				Name:       podmetrics.Name,
				Namespace:  podmetrics.Namespace,
				Containers: []msg.ContainerResource{},
			}
			for _, containermetrics := range podmetrics.Containers {
				celem := msg.ContainerResource{
					Name:   containermetrics.Name,
					CpuUse: containermetrics.Usage.Cpu().MilliValue(),
					MemUse: containermetrics.Usage.Memory().MilliValue(),
				}
				pelem.Containers = append(pelem.Containers, celem)
			}
			res.Pods = append(res.Pods, pelem)
		}
	}

	return res, nil
}
