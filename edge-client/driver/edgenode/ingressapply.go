package edgenode

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// IngressApplyByConfig Ingress部署
func (ed *Driver) IngressApplyByConfig(ingressConfig *v1beta1.Ingress, namespace string) error {

	ingressClient := ed.clientset.ExtensionsV1beta1().Ingresses(namespace)
	if _, err := ingressClient.Get(context.TODO(), ingressConfig.Name, metav1.GetOptions{}); err == nil {
		// 更新
		ingressApi, err := ingressClient.Update(context.TODO(), ingressConfig, metav1.UpdateOptions{})
		if err != nil {
			// err = fmt.Errorf("podClient.Update error {%v}", err)
			log.Error(err)
			return err
		}
		log.Debugf("namespace {%v} ingress {%v} update success", ingressApi.Namespace, ingressApi.Name)
	} else {
		// 创建
		ingressApi, err := ingressClient.Create(context.TODO(), ingressConfig, metav1.CreateOptions{})
		if err != nil {
			// err = fmt.Errorf("podClient.Update error {%v}", err)
			log.Error(err)
			return err
		}
		log.Debugf("namespace {%v} svc {%v} create success", ingressApi.Namespace, ingressApi.Name)
	}
	return nil
}

// IngressApply Ingress部署
func (ed *Driver) IngressApply(ingressYaml []byte, namespace string) error {

	ingressJson, err := yaml.ToJSON(ingressYaml)
	if err != nil {
		err = fmt.Errorf("yaml.ToJSON error {%v}", err)
		log.Error(err)
		return err
	}

	ingressConfig := &v1beta1.Ingress{}
	if err := json.Unmarshal(ingressJson, ingressConfig); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return err
	}

	return ed.IngressApplyByConfig(ingressConfig, namespace)
}
