package api

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// PodYamlDownload Pod配置文件下载
func (ad *Driver) PodYamlDownload(filehash string) ([]byte, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", ad.token).
		SetQueryParams(map[string]string{
			"filehash": filehash,
		}).Get(ad.info.PodYamlDownload.URL)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return nil, err
	}

	isJSON := false
	isBinary := false
	for _, elem := range resp.Header()["Content-Type"] {
		if strings.Contains(elem, "application/json") {
			isJSON = true
		} else if strings.Contains(elem, "application/octet-stream") {
			isBinary = true
		}
	}

	if isJSON {
		var result response.Response
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			log.Errorf("json Unmarshal error {%v}", err)
			return nil, err
		}
		err := fmt.Errorf("file download error {%v}", result.Msg)
		log.Error(err)
		return nil, err
	} else if isBinary {
		return resp.Body(), nil
	} else {
		err := fmt.Errorf("http Content-Type error {%v}", resp.Header()["Content-Type"])
		log.Error(err)
		return nil, err
	}
}
