package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"

	log "github.com/sirupsen/logrus"
)

// GetClientIDByUserName 利用username获取客户端的id
func GetClientIDByUserName(username string) (string, error) {
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		log.Errorf("redis get rdb error {%s}", err)
		return "", err
	}

	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("clients:id:%s:token:*", username)).Result()
	if err != nil {
		log.Errorf("redis keys error {%s}", err)
		return "", err
	}
	if len(tkeys) < 1 || len(tkeys) > 1 {
		log.Errorf("redis keys length error {%v}", tkeys)
		return "", err
	}

	valuestr, err := rdb.Get(context.TODO(), tkeys[0]).Result()
	if err != nil {
		err = fmt.Errorf("redis get {%v} error {%v}", tkeys[0], err)
		log.Error(err)
		return "", err
	}

	clientloginstatus := value.ClientLoginStatus{}
	if err := json.Unmarshal([]byte(valuestr), &clientloginstatus); err != nil {
		err = fmt.Errorf("json.Unmarshal {%v} error {%v}", valuestr, err)
		log.Error(err)
		return "", err
	}

	return clientloginstatus.ClientID, nil
}
