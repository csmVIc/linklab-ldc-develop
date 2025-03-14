package edgenode

import (
	"errors"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Driver 负责边缘设备操作
type Driver struct {
	info      *EInfo
	k8sconfig *rest.Config
	clientset *kubernetes.Clientset
	edgenodes EdgeNodes
	pods      Pods
}

// EDriver 全局边缘设备操作实例
var (
	EDriver *Driver
)

func (ed *Driver) init() error {

	var err error = nil
	if ed.info.K8SClient.InCluster {
		// K8S 集群内初始化
		ed.k8sconfig, err = rest.InClusterConfig()
		if err != nil {
			err = fmt.Errorf("rest.InClusterConfig error {%v}", err)
			log.Error(err)
			return err
		}
	} else {
		// K8S 集群外初始化
		ed.k8sconfig, err = clientcmd.BuildConfigFromFlags("", ed.info.K8SClient.KubeConfig)
		if err != nil {
			err = fmt.Errorf("clientcmd.BuildConfigFromFlags error {%v}", err)
			log.Error(err)
			return err
		}
	}

	ed.clientset, err = kubernetes.NewForConfig(ed.k8sconfig)
	if err != nil {
		err = fmt.Errorf("kubernetes.NewForConfig error {%v}", err)
		log.Error(err)
		return err
	}
	return nil
}

// New 创建设备
func New(i *EInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info i nil error")
	}
	ed := &Driver{info: i, edgenodes: EdgeNodes{
		edgeNodesMap:  sync.Map{},
		edgeNodesLock: sync.RWMutex{},
	}, pods: Pods{
		podsMap:  sync.Map{},
		podsLock: sync.RWMutex{},
	}}
	if err := ed.init(); err != nil {
		return nil, fmt.Errorf("init error {%v}", err)
	}
	return ed, nil
}
