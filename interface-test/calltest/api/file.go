package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/parameter/response"
	"math/rand"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// FileUpload 烧写文件上传
func (ad *Driver) FileUpload(boardname string, filesize int, filerandom bool, token string, filepath string, compileoutput []byte) (string, error) {

	parameter := map[string]string{
		"boardname": boardname,
	}
	pbuffer, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		log.Error(err)
		return "", err
	}

	binary := make([]byte, filesize)

	// 随机化文件二进制值
	if len(compileoutput) > 0 {
		binary = compileoutput
	} else if filerandom {
		_, err := rand.Read(binary)
		if err != nil {
			err = fmt.Errorf("rand.Read error {%v}", err)
			log.Error(err)
			return "", err
		}
	} else {
		binary, err = ioutil.ReadFile(filepath)
		if err != nil {
			err = fmt.Errorf("ioutil.ReadFile error {%v}", err)
			log.Error(err)
			return "", err
		}
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", token).
		SetMultipartFields(
			&resty.MultipartField{
				Param:       "parameters",
				ContentType: "application/json",
				Reader:      bytes.NewReader(pbuffer),
			},
			&resty.MultipartField{
				Param:       "file",
				FileName:    "burn.bin",
				ContentType: "application/octet-stream",
				Reader:      bytes.NewReader(binary),
			},
		).
		Post(ad.info.FileUpload.URL)
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
		err := fmt.Errorf("file upload error {%v}", result.Msg)
		log.Error(err)
		return "", err
	}

	fileHashMap := result.Data.(map[string]interface{})
	if _, isOk := fileHashMap["filehash"]; isOk == false {
		err = fmt.Errorf("payload not contain {filehash} {%v}", result)
		log.Error(err)
		return "", err
	}

	return fileHashMap["filehash"].(string), nil
}
