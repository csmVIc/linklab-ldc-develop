package edgenode

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// labelNode 增加标签
func (ed *Driver) labelNode(nodename string, key string, value string) error {

	labelPatch := []byte(fmt.Sprintf(`[{"op":"add","path":"/metadata/labels/%s","value":"%s" }]`, key, value))
	if _, err := ed.clientset.CoreV1().Nodes().Patch(context.TODO(), nodename, types.JSONPatchType, labelPatch, metav1.PatchOptions{}); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// unlabelNode 删除标签
func (ed *Driver) unlabelNode(nodename string, key string) error {

	unlabelPatch := []byte(fmt.Sprintf(`[{"op":"remove","path":"/metadata/labels/%s"}]`, key))
	if _, err := ed.clientset.CoreV1().Nodes().Patch(context.TODO(), nodename, types.JSONPatchType, unlabelPatch, metav1.PatchOptions{}); err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// labelEdgeType 标记边缘类型 
func (ed *Driver) labelEdgeType(nodename string) error {

	edgetype := nodename[:strings.LastIndex(nodename, "-")]
	if strings.ToLower(edgetype) == "jetsonnano"{
		// 对于JetsonNano, 需要特殊标记NodeName
		if err := ed.labelNode(nodename, "linklab.edgetype", nodename); err != nil {
			err = fmt.Errorf("label edgenode {%v} nodename {%v} error {%v}", nodename, nodename, err)
			log.Error(err)
			return err
		}
	}else{
		if err := ed.labelNode(nodename, "linklab.edgetype", edgetype); err != nil {
			err = fmt.Errorf("label edgenode {%v} edgetype {%v} error {%v}", nodename, edgetype, err)
			log.Error(err)
			return err
		}
	}

	return nil
}

