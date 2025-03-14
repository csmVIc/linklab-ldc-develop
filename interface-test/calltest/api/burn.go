package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// DeviceBurn 设备烧写
func (ad *Driver) DeviceBurn(parameter *request.DeviceBurnTasks, token string) (string, error) {

	if parameter == nil {
		err := errors.New("parameter nil error")
		log.Error(err)
		return "", err
	}

	body, err := json.Marshal(*parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		log.Error(err)
		return "", err
	}

	client := resty.New().SetTimeout(time.Minute * 2)
	resp, err := client.R().
		SetHeader("Authorization", token).
		SetBody(body).
		Post(ad.info.Burn.URL)
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
		err := fmt.Errorf("device burn error {%v}", result.Msg)
		log.Error(err)
		return "", err
	}

	groupIDMap := result.Data.(map[string]interface{})
	if _, isOk := groupIDMap["groupid"]; isOk == false {
		err = fmt.Errorf("payload not contain {groupid} {%v}", result)
		log.Error(err)
		return "", err
	}

	return groupIDMap["groupid"].(string), nil
}
