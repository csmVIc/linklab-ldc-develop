package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// UploadEdgePodYaml 上传边缘Pod配置文件
func (ud *Driver) UploadEdgePodYaml(token string, yamlbinary []byte) (string, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", token).
		SetMultipartFields(
			&resty.MultipartField{
				Param:       "file",
				FileName:    "pod.yaml",
				ContentType: "application/octet-stream",
				Reader:      bytes.NewReader(yamlbinary),
			},
		).Post(ud.info.EdgeNode.PodYamlUploadURL)
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
		err := fmt.Errorf("upload pod yaml error {%v}", result.Msg)
		// log.Error(err)
		return "", err
	}

	if _, isOk := (*datamap)["filehash"]; !isOk {
		err := fmt.Errorf("response doesn't contain filehash error")
		// log.Error(err)
		return "", err
	}

	return (*datamap)["filehash"], nil
}
