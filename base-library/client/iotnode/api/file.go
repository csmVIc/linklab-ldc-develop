package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// FileDownload 文件下载
func (ad *Driver) FileDownload(boardname string, filehash string, groupid string, taskindex int, suffix string) (string, error) {

	// 重新生成文件存储哈希值
	regenhash := ad.regenFileHash(boardname, filehash, groupid, taskindex)

	// 检查该文件是否已经下载
	filePath := fmt.Sprintf("%v/%v.%v", ad.info.TmpDir, regenhash, suffix)
	if _, err := os.Stat(filePath); os.IsExist(err) {
		return filePath, nil
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", ad.token).
		SetQueryParams(map[string]string{
			"filehash":  filehash,
			"boardname": boardname,
		}).Get(ad.info.FileDownload.URL)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return "", err
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
			return "", err
		}
		err := fmt.Errorf("file download error {%v}", result.Msg)
		log.Error(err)
		return "", err
	} else if isBinary {
		if err := ioutil.WriteFile(filePath, resp.Body(), 0644); err != nil {
			log.Errorf("file write error {%v}", err)
			return "", err
		}
		return filePath, nil
	} else {
		err := fmt.Errorf("http Content-Type error {%v}", resp.Header()["Content-Type"])
		log.Error(err)
		return "", err
	}
}
