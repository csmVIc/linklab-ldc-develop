package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// GetLoginStatus 获取登录状态
func (cd *Driver) GetLoginStatus(utype string, id string, loginstatus interface{}) error {

	if cd.rdb == nil {
		err := errors.New("cd.rdb is nil error")
		log.Error(err)
		return err
	}

	if err := cd.ping(); err != nil {
		err := fmt.Errorf("cd.ping error {%v}", err)
		log.Error(err)
		return err
	}

	tkeys, err := cd.rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:%s:token:*", utype, id)).Result()
	if err != nil {
		err := fmt.Errorf("cd.rdb keys error {%v}", err)
		log.Error(err)
		return err
	}

	if len(tkeys) < 1 || len(tkeys) > 1 {
		err := fmt.Errorf("keys len {%v} error", len(tkeys))
		log.Error(err)
		return err
	}

	valuestr, err := cd.rdb.Get(context.TODO(), tkeys[0]).Result()
	if err != nil {
		err := fmt.Errorf("cd.rdb get {%v} error {%v}", tkeys[0], err)
		log.Error(err)
		return err
	}

	if err := json.Unmarshal([]byte(valuestr), loginstatus); err != nil {
		err := fmt.Errorf("json.Unmarshal error {%v}", err)
		log.Error(err)
		return err
	}

	return nil
}
