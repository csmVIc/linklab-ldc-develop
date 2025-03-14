package edgenode

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// ServiceApplyByConfig 服务部署
func (ed *Driver) ServiceApplyByConfig(svcConfig *corev1.Service, namespace string) error {

	svcClient := ed.clientset.CoreV1().Services(namespace)
	if _, err := svcClient.Get(context.TODO(), svcConfig.Name, metav1.GetOptions{}); err == nil {
		// 更新服务
		svcApi, err := svcClient.Update(context.TODO(), svcConfig, metav1.UpdateOptions{})
		if err != nil {
			// err = fmt.Errorf("podClient.Update error {%v}", err)
			log.Error(err)
			return err
		}
		log.Debugf("namespace {%v} svc {%v} update success", svcApi.Namespace, svcApi.Name)
	} else {
		// 创建服务
		svcApi, err := svcClient.Create(context.TODO(), svcConfig, metav1.CreateOptions{})
		if err != nil {
			// err = fmt.Errorf("podClient.Update error {%v}", err)
			log.Error(err)
			return err
		}
		log.Debugf("namespace {%v} svc {%v} create success", svcApi.Namespace, svcApi.Name)
	}
	return nil
}

// ServiceApply 服务部署
func (ed *Driver) ServiceApply(serviceYaml []byte, namespace string) error {

	svcJson, err := yaml.ToJSON(serviceYaml)
	if err != nil {
		err = fmt.Errorf("yaml.ToJSON error {%v}", err)
		log.Error(err)
		return err
	}

	svcConfig := &corev1.Service{}
	if err := json.Unmarshal(svcJson, svcConfig); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return err
	}

	return ed.ServiceApplyByConfig(svcConfig, namespace)
}
