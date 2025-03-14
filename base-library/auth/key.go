package auth

import (
	"context"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"strings"

	log "github.com/sirupsen/logrus"
)

// GetIDFromKey 从键值中解析出id值
func GetIDFromKey(key string) (string, error) {
	arr := strings.Split(key, ":")
	if len(arr) < 3 {
		err := fmt.Errorf("key {%v} split len < 3", key)
		log.Error(err)
		return "", err
	}
	return arr[2], nil
}

// GetTokenFromKey 从键值中解析出token值
func GetTokenFromKey(key string) (string, error) {
	arr := strings.Split(key, ":")
	if len(arr) < 5 {
		err := fmt.Errorf("key {%v} split len < 5", key)
		log.Error(err)
		return "", err
	}
	return arr[4], nil
}

// GetTokenFromID 通过ID获得token值
func GetTokenFromID(utype string, id string) (string, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		log.Errorf("redis get rdb error {%s}", err)
		return "", err
	}

	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:%s:token:*", utype, id)).Result()
	if err != nil {
		log.Errorf("redis keys error {%s}", err)
		return "", err
	}

	if len(tkeys) < 1 || len(tkeys) > 1 {
		log.Errorf("redis keys length error {%v}", tkeys)
		return "", err
	}

	token, err := GetTokenFromKey(tkeys[0])
	if err != nil {
		err = fmt.Errorf("GetTokenFromKey error {%v}", err)
		log.Error(err)
		return "", err
	}

	return token, nil
}
