package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ApplyEdgePod 部署边缘Pod
func (ud *Driver) ApplyEdgePod(token string, yamlhash string, edgeregistry bool, createingress bool) (string, map[string]string, error) {

	parameter := request.UserPodApply{
		YamlHash:        yamlhash,
		UseEdgeRegistry: edgeregistry,
		CreateIngress:   createingress,
	}

	body, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		// log.Error(err)
		return "", nil, err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post(ud.info.EdgeNode.PodApplyURL)
	if err != nil {
		err = fmt.Errorf("http post error {%v}", err)
		// log.Error(err)
		return "", nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		// log.Error(err)
		return "", nil, err
	}

	respData := &response.PodApply{}
	result := response.Response{
		Data: respData,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		// log.Error(err)
		return "", nil, err
	}

	if result.Code != 0 {
		err := fmt.Errorf("%v", result.Msg)
		// log.Error(err)
		return "", nil, err
	}

	return respData.ClientID, respData.IngressMap, nil
}
