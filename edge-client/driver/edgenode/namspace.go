package edgenode

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceCreateIfNotExist 如果namespace不存在，那么创建
func (ed *Driver) NamespaceCreateIfNotExist(namespace string) error {
	namespaceClient := ed.clientset.CoreV1().Namespaces()
	if _, err := namespaceClient.Get(context.TODO(), namespace, metav1.GetOptions{}); err == nil {
		// namespace已存在
		return nil
	}

	namespaceConfig := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	if _, err := namespaceClient.Create(context.TODO(), namespaceConfig, metav1.CreateOptions{}); err != nil {
		err := fmt.Errorf("namespace {%v} create error {%v}", namespace, err)
		log.Error(err)
		return err
	}
	log.Debugf("namespace {%v} create success", namespace)
	return nil
}

// lsNamespace 显示所有Namespace
func (ed *Driver) lsNamespace() ([]string, error) {
	namespaceClient := ed.clientset.CoreV1().Namespaces()
	namespaceList, err := namespaceClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		err := fmt.Errorf("namespace list error {%v}", err)
		log.Error(err)
		return nil, err
	}

	res := []string{}
	for _, ns := range namespaceList.Items {
		if _, isSystem := ed.info.K8SPod.SystemNamespaces[ns.Name]; isSystem == true {
			// 系统命名空间直接跳过
			continue
		}
		res = append(res, ns.Name)
	}

	return res, nil
}

// deleteNamespaceIfPodEmpty 如果命名空间下的Pod为空，删除命名空间
func (ed *Driver) deleteNamespaceIfPodEmpty(namespace string) error {

	if ed.info.K8SPod.SystemNamespaces[namespace] == true || namespace == "default" {
		return nil
	}

	podClient := ed.clientset.CoreV1().Pods(namespace)
	podList, err := podClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("podClient.List error {%v}", err)
		log.Error(err)
		return err
	}

	if len(podList.Items) > 0 {
		return nil
	}

	deletePolicy := metav1.DeletePropagationForeground
	namespaceClient := ed.clientset.CoreV1().Namespaces()
	if err := namespaceClient.Delete(context.TODO(), namespace, metav1.DeleteOptions{PropagationPolicy: &deletePolicy}); err != nil {
		err = fmt.Errorf("namespaceClient.Delete error {%v}", err)
		log.Error(err)
		return err
	}

	log.Debugf("namespace delete {%v} success", namespace)

	return nil
}
