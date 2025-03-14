package edgenode

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IngressDelete Ingress删除
func (ed *Driver) IngressDelete(namespace string, name string) error {

	ingressClient := ed.clientset.ExtensionsV1beta1().Ingresses(namespace)

	if _, err := ingressClient.Get(context.TODO(), name, metav1.GetOptions{}); err != nil {
		log.Error(err)
		return nil
	}

	deletePolicy := metav1.DeletePropagationForeground
	if err := ingressClient.Delete(context.TODO(), name, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		// err = fmt.Errorf("podClient.Delete error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}
