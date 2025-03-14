package user

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ListEdgePod 显示边缘服务
func (ud *Driver) ListEdgePod(token string) (*[]response.PodInfo, error) {

	client := resty.New()
	resp, err := client.R().
		SetHeader("Authorization", token).
		Get(ud.info.EdgeNode.PodListURL)
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

	podlist := &response.PodList{}
	result := response.Response{
		Data: podlist,
	}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		err := fmt.Errorf("json Unmarshal error {%v}", err)
		// log.Error(err)
		return nil, err
	}

	if result.Code != 0 {
		err := fmt.Errorf("list edge pod error {%v}", result.Msg)
		// log.Error(err)
		return nil, err
	}

	return &podlist.Pods, nil
}
