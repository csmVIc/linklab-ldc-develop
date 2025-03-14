package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// BuildEdgeImage 构建边缘镜像
func (ud *Driver) BuildEdgeImage(token string, filehash string, imagename string, nodeselector map[string]string) (string, error) {

	parameter := map[string]interface{}{
		"filehash":     filehash,
		"imagename":    imagename,
		"nodeselector": nodeselector,
	}

	body, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		// log.Error(err)
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post(ud.info.EdgeNode.ImageBuildURL)
	if err != nil {
		err = fmt.Errorf("http post error {%v}", err)
		// log.Error(err)
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		// log.Error(err)
		return "", err
	}

	datamap := &map[string]string{}
	result := response.Response{
		Data: datamap,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		// log.Error(err)
		return "", err
	}

	if result.Code != 0 {
		err := fmt.Errorf("%v", result.Msg)
		// log.Error(err)
		return "", err
	}

	if _, isOk := (*datamap)["clientid"]; !isOk {
		err := fmt.Errorf("response doesn't contain clientid error")
		// log.Error(err)
		return "", err
	}

	return (*datamap)["clientid"], nil
}
