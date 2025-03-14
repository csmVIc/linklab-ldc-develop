package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// CreateDeviceBindGroup 创建设备绑定组
func (ud *Driver) CreateDeviceBindGroup(token string, grouptype string, boardtypes []string) error {

	parameter := map[string]interface{}{
		"type":   grouptype,
		"boards": boardtypes,
	}

	body, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		// log.Error(err)
		return err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post(ud.info.DeviceBindGroup.CreateGroupURL)
	if err != nil {
		err = fmt.Errorf("http post error {%v}", err)
		// log.Error(err)
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		// log.Error(err)
		return err
	}

	result := response.Response{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		// log.Error(err)
		return err
	}

	if result.Code != 0 {
		err := fmt.Errorf("create device bind group error {%v}", result.Msg)
		// log.Error(err)
		return err
	}

	return nil
}
