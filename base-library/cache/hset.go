package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// HSetAndExpire 哈希表设置并且设置超时时间
func (cd *Driver) HSetAndExpire(key string, values interface{}, expiration time.Duration) error {
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
	if _, err := cd.rdb.HSet(context.TODO(), key, values).Result(); err != nil {
		err = fmt.Errorf("redis rdb hset error {%s}", err)
		log.Error(err)
		return err
	}
	if _, err := cd.rdb.Expire(context.TODO(), key, expiration).Result(); err != nil {
		err = fmt.Errorf("redis rdb expire error {%s}", err)
		log.Error(err)
		return err
	}
	return nil
}
