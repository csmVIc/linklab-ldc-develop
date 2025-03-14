package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"

	log "github.com/sirupsen/logrus"
)

// checkLoginStatus 检查登录状态，如果没有问题返回id
func (ah *Handler) checkLoginStatus(token string) (string, error) {
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return "", err
	}
	log.Debugf("成功获取 Redis 客户端: %+v", rdb) // 调试日志，输出 Redis 客户端信息

	keys, err := rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:*:token:%s", ah.uType, token)).Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%s}", err)
		log.Error(err)
		return "", err
	}
	log.Debugf("从 Redis 获取到的 keys: %+v", keys) // 调试日志，输出从 Redis 获取的 keys

	if len(keys) < 1 || len(keys) > 1 {
		err = fmt.Errorf("redis keys length {%v} error", len(keys))
		log.Error(err)
		return "", err
	}
	log.Debugf("找到有效的 key: %s", keys[0]) // 调试日志，输出有效的 key

	id, err := GetIDFromKey(keys[0])
	if err != nil {
		err = fmt.Errorf("get id from key error {%v}", err)
		log.Error(err)
		return "", err
	}
	log.Debugf("成功获取 ID: %s", id) // 调试日志，输出获取到的 ID

	return id, nil
}

// CheckClientIDAndGetToken 检查客户端的登录状态
func CheckClientIDAndGetToken(username string, clientid string) (string, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return "", err
	}

	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("clients:id:%s:token:*", username)).Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%s}", err)
		log.Error(err)
		return "", err
	}
	if len(tkeys) < 1 || len(tkeys) > 1 {
		err = fmt.Errorf("redis keys length error {%v}", tkeys)
		log.Error(err)
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

	if clientloginstatus.ClientID != clientid {
		err = fmt.Errorf("client {%v} token clientid {%v} != {%v} error", username, clientloginstatus.ClientID, clientid)
		log.Error(err)
		return "", err
	}

	return tkeys[0], err
}
