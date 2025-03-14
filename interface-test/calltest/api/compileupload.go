package api

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// CompileUpload 编译上传
func (ad *Driver) CompileUpload(boardtype string, compiletype string, filepath string) (string, error) {

	binary, err := ioutil.ReadFile(filepath)
	if err != nil {
		err = fmt.Errorf("ioutil.ReadFile error {%v}", err)
		log.Error(err)
		return "", err
	}

	filehash := fmt.Sprintf("%x", sha1.Sum(binary))
	parameter := map[string]string{
		"filehash":    filehash,
		"boardType":   boardtype,
		"compileType": compiletype,
	}
	pbuffer, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		log.Error(err)
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		SetMultipartFields(
			&resty.MultipartField{
				Param:       "parameters",
				ContentType: "application/json",
				Reader:      bytes.NewReader(pbuffer),
			},
			&resty.MultipartField{
				Param:       "file",
				FileName:    "compile.zip",
				ContentType: "application/octet-stream",
				Reader:      bytes.NewReader(binary),
			},
		).
		Post(ad.info.CompileUpload.URL)
	if err != nil {
		err = fmt.Errorf("http post error {%v}", err)
		log.Error(err)
		return "", err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return "", err
	}

	result := response.Response{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		log.Error(err)
		return "", err
	}

	if result.Code != 0 {
		err := fmt.Errorf("compile upload error {%v}", result.Msg)
		log.Error(err)
		return "", err
	}

	log.Infof("compile upload success {%v}", result.Msg)
	return filehash, nil
}
