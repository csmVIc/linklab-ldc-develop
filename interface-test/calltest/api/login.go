package api

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// UserLogin 用户登录
func (ad *Driver) UserLogin(username string, password string) (string, error) {

	parameter := map[string]string{
		"id":       username,
		"password": password,
	}

	body, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		log.Error(err)
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(ad.info.Login.URL)
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
		err := fmt.Errorf("user login error {%v}", result.Msg)
		log.Error(err)
		return "", err
	}

	tokenMap := result.Data.(map[string]interface{})
	if _, isOk := tokenMap["token"]; isOk == false {
		err = fmt.Errorf("payload not contain {token} {%v}", result)
		log.Error(err)
		return "", err
	}

	return tokenMap["token"].(string), nil
}
