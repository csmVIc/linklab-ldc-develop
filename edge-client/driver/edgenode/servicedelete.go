package edgenode

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceDelete 服务删除
func (ed *Driver) ServiceDelete(namespace string, service string) error {

	svcClient := ed.clientset.CoreV1().Services(namespace)

	if _, err := svcClient.Get(context.TODO(), service, metav1.GetOptions{}); err != nil {
		log.Error(err)
		return nil
	}

	deletePolicy := metav1.DeletePropagationForeground
	if err := svcClient.Delete(context.TODO(), service, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		// err = fmt.Errorf("podClient.Delete error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}
