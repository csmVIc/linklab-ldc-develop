package saasbackend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

// GetUserSiteInfo 获取用户的Site信息
func (sd *Driver) GetUserSiteInfo(userId string) (*SiteResult, error) {

	client := resty.New()
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"uid": userId,
		}).Get(sd.info.UserSite.URL)

	if err != nil {
		log.Errorf("http error {%v}", err)
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		err := fmt.Errorf("http status code error {%v}", resp.StatusCode())
		log.Error(err)
		return nil, err
	}

	result := SaasResponse{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Errorf("json Unmarshal error {%v}", err)
		return nil, err
	}

	if result.Code != 0 {
		err = fmt.Errorf("result.Code {%v} != 0 error {%v}", result.Code, result.Message)
		log.Error(err)
		return nil, err
	}

	data := result.Data.(map[string]interface{})
	siteResult := &SiteResult{
		ID:   int(data["id"].(float64)),
		Name: data["siteName"].(string),
	}

	log.Debugf("user {%v} get site info id {%v} name {%v}", userId, siteResult.ID, siteResult.Name)
	return siteResult, nil
}
