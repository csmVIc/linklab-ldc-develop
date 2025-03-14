package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"strings"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// CompileDownload 编译下载
func (ad *Driver) CompileDownload(boardtype string, compiletype string, filehash string) ([]byte, bool, error) {

	client := resty.New()
	resp, err := client.R().SetQueryParams(map[string]string{
		"filehash":    filehash,
		"boardtype":   boardtype,
		"compiletype": compiletype,
	}).Get(ad.info.CompileDownload.URL)
	if err != nil {
		err = fmt.Errorf("http get error {%v}", err)
		log.Error(err)
		return []byte{}, false, err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return []byte{}, false, err
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
		result := response.Response{}
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			log.Error(err.Error())
			return []byte{}, false, err
		}

		if result.Code != 0 {
			err := fmt.Errorf("compile download error {%v}", result.Msg)
			log.Error(err)
			return []byte{}, false, err
		}

		if result.Msg == "error" {
			err := fmt.Errorf("compile download builderr error {%v}", result.Msg)
			log.Error(err)
			return []byte{}, false, err
		}

		if result.Msg == "notexist" {
			err := fmt.Errorf("compile download notexist error {%v}", result.Msg)
			log.Error(err)
			return []byte{}, true, err
		}

		log.Debugf("compile status {%v}", result.Msg)
		return []byte{}, false, nil
	} else if isBinary {
		log.Debugf("compile download success")
		return resp.Body(), false, nil
	} else {
		err := errors.New("http Content-Type error")
		log.Errorf("compile download content-type error {%v}", resp.Header())
		return []byte{}, false, err
	}
}
