package sub

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (sd *Driver) getPodInfo(yamlhash string) (string, map[string]string, error) {

	// 读取文件
	filter := &table.PodYamlFilter{
		FileHash: yamlhash,
	}
	result := &table.PodYaml{}
	if err := database.Mdriver.FindOneElem("podyaml", filter, result); err != nil {
		// 需要的文件不存在
		err = fmt.Errorf("podyaml filehash {%v} not exist error {%v}", yamlhash, err)
		log.Error(err)
		return "", nil, err
	}

	podJson, err := yaml.ToJSON(result.FileData.Data)
	if err != nil {
		err = fmt.Errorf("yaml.ToJSON error {%v}", err)
		log.Error(err)
		return "", nil, err
	}

	podConfig := &apiv1.Pod{}
	if err := json.Unmarshal(podJson, podConfig); err != nil {
		err = fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return "", nil, err
	}

	return podConfig.Name, podConfig.Spec.NodeSelector, nil
}
