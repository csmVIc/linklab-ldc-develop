package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// LinkDeviceBindGroup 关联设备绑定组
func (ud *Driver) LinkDeviceBindGroup(token string, grouptype string, devices []request.DevInfoForGroup) (string, error) {

	parameter := map[string]interface{}{
		"type":    grouptype,
		"devices": devices,
	}

	body, err := json.Marshal(parameter)
	if err != nil {
		err = fmt.Errorf("json Marshal error {%v}", err)
		// log.Error(err)
		return "", err
	}

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", token).
		SetBody(body).
		Post(ud.info.DeviceBindGroup.LinkGroupURL)
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
		err := fmt.Errorf("%v", result.Msg)
		// log.Error(err)
		return "", err
	}

	if _, isOk := (*datamap)["id"]; !isOk {
		err := fmt.Errorf("response doesn't contain id error")
		// log.Error(err)
		return "", err
	}

	return (*datamap)["id"], nil
}
