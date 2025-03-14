package driver

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// buildDownload 构建文件下载
func (dd *Driver) buildDownload() ([]byte, error) {
// 使用 resty 创建了一个 HTTP 客户端 client，用于发起 HTTP 请求。
	client := resty.New()
	//  发起 HTTP 请求
	resp, err := client.R().
		SetHeader("Authorization", dd.info.API.Token).
		SetQueryParams(map[string]string{
			"filehash": dd.info.API.FileHash,
		}).Get(dd.info.API.BuildDownloadURL)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return nil, err
	}
// 如果服务器返回的状态码不是 200 OK，记录错误日志并返回错误。
	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return nil, err
	}

	isJSON := false
	isBinary := false
	// 遍历响应头中的 Content-Type，判断文件类型
	// _ 表示忽略索引。	elem 是当前遍历的值（string 类型），例如 "application/json"。
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
