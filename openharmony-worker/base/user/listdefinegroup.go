package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ListDefineDeviceBindGroup 显示已定义设备绑定组
func (ud *Driver) ListDefineDeviceBindGroup(token string) (*[]response.BindGroupInfo, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", token).
		Get(ud.info.DeviceBindGroup.ListDefineGroupURL)
	if err != nil {
		err = fmt.Errorf("http get error {%v}", err)
		// log.Error(err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		// log.Error(err)
		return nil, err
	}

	grouplist := &response.BindGroupInfoList{}
	result := response.Response{
		Data: grouplist,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		// log.Error(err)
		return nil, err
	}

	if result.Code != 0 {
		err := fmt.Errorf("list define device bind group error {%v}", result.Msg)
		// log.Error(err)
		return nil, err
	}

	return &grouplist.Groups, nil
}
