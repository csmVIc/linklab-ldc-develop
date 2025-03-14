package login

import (
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"

	log "github.com/sirupsen/logrus"
)

func getTrueCheckAndSaltFromValue(valuestr string) (string, string, error) {

	loginstatus := value.UserLoginStatus{}
	if err := json.Unmarshal([]byte(valuestr), &loginstatus); err != nil {
		err := fmt.Errorf("json.Unmarshal {%v} error {%v}", valuestr, err)
		log.Error(err)
		return "", "", err
	}

	return loginstatus.TrueCheck, loginstatus.Salt, nil
}

func getTrueCheckAndSaltAndTenantIDFromDb(id string, utype string) (string, string, int, error) {

	trueCheck, salt, tenantid := "", "", -1
	switch utype {
	case "users":
		user := &table.User{}
		filter := table.UserFilter{UserID: id}
		err := database.Mdriver.FindOneElem(utype, filter, user)
		if err != nil {
			err = fmt.Errorf("{%v} database.Mdriver.FindOneElem error {%v}", id, err)
			log.Error(err)
			return "", "", -1, err
		}
		trueCheck = user.Hash
		salt = user.Salt
		tenantid = user.TenantID
	case "clients":
		client := &table.Client{}
		filter := table.ClientFilter{UserName: id}
		err := database.Mdriver.FindOneElem(utype, filter, client)
		if err != nil {
			err = fmt.Errorf("{%v} database.Mdriver.FindOneElem error {%v}", id, err)
			log.Error(err)
			return "", "", -1, err
		}
		trueCheck = client.Password
		salt = client.Salt
	default:
		err := fmt.Errorf("unsupported user type {%v} error", utype)
		log.Error(err)
		return "", "", -1, err
	}

	return trueCheck, salt, tenantid, nil
}
