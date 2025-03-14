package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ListLinkDeviceBindGroup 显示已关联设备绑定组
func (ud *Driver) ListLinkDeviceBindGroup(token string) (*[]response.DeviceGroupInfo, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", token).
		Get(ud.info.DeviceBindGroup.ListLinkGroupURL)
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

	grouplist := &response.DeviceGroupInfoList{}
	result := response.Response{
		Data: grouplist,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		// log.Error(err)
		return nil, err
	}

	if result.Code != 0 {
		err := fmt.Errorf("list link device bind group error {%v}", result.Msg)
		// log.Error(err)
		return nil, err
	}

	return &grouplist.Groups, nil
}
