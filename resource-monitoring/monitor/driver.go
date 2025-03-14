package monitor

import (
	"errors"
	"sync"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Driver 资源监控
type Driver struct {
	info      *MInfo
	exitwg    sync.WaitGroup
	clientset *kubernetes.Clientset
}

// Md 资源监控全局实例
var (
	Mdriver *Driver
)

// k8s客户端连接初始化
func (md *Driver) k8sinit() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("rest.InClusterConfig error {%v}", err)
		return err
	}
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("kubernetes.NewForConfig error {%v}", err)
		return err
	}
	md.clientset = cs
	return nil
}

// New 创建资源监控实例
func New(i *MInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info is nil error")
	}
	md := &Driver{info: i, exitwg: sync.WaitGroup{}, clientset: nil}
	err := md.k8sinit()
	if err != nil {
		return md, err
	}
	return md, nil
}
