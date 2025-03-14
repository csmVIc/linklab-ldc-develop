package edgenode

import (
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodDelete Pod删除
func (ed *Driver) PodDelete(namespace string, pod string) error {
	podClient := ed.clientset.CoreV1().Pods(namespace)

	deletePolicy := metav1.DeletePropagationForeground
	if err := podClient.Delete(context.TODO(), pod, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		// err = fmt.Errorf("podClient.Delete error {%v}", err)
		log.Error(err)
		return err
	}

	// 启动命名空间回收
	go func() {
		for {
			time.Sleep(time.Second)
			_, err := podClient.Get(context.TODO(), pod, metav1.GetOptions{})
			if err != nil {
				break
			}
		}

		log.Debugf("namespace {%v} pod {%v} delete success", namespace, pod)

		if err := ed.deleteNamespaceIfPodEmpty(namespace); err != nil {
			err = fmt.Errorf("ed.deleteNamespaceIfPodEmpty error {%v}", err)
			log.Error(err)
			return
		}
	}()

	return nil
}
