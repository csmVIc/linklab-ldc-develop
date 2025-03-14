package cache

import (
	"fmt"
	"linklab/device-control-v2/base-library/cache/value"

	log "github.com/sirupsen/logrus"
)

// GetUserTenantID 获取用户的租户ID
func (cd *Driver) GetUserTenantID(userid string) (int, error) {

	loginstatus := value.UserLoginStatus{}
	if err := cd.GetLoginStatus("users", userid, &loginstatus); err != nil {
		err := fmt.Errorf("cd.GetLoginStatus {users} {%v} error {%v}", userid, err)
		log.Error(err)
		return -1, err
	}

	return loginstatus.TenantID, nil
}

// GetClientTenantID 获取客户端的租户ID
func (cd *Driver) GetClientTenantID(clientid string) (map[int]bool, error) {

	loginstatus := value.ClientLoginStatus{}
	if err := cd.GetLoginStatus("clients", clientid, &loginstatus); err != nil {
		err := fmt.Errorf("cd.GetLoginStatus {clients} {%v} error {%v}", clientid, err)
		log.Error(err)
		return nil, err
	}
	return loginstatus.TenantID, nil
}
