package cache

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func (cd *Driver) lockToken(lockid string) string {
	ntime := time.Now().Unix()
	rint := rand.Int()
	tokenseed := fmt.Sprintf("%s:%v:%v", lockid, ntime, rint)
	tokenbinary := sha256.Sum256([]byte(tokenseed))
	return fmt.Sprintf("%x", tokenbinary)
}

// Lock 分布式锁加锁
func (cd *Driver) Lock(lockid string) (string, error) {
	if cd.rdb == nil {
		err := errors.New("cd.rdb is nil error")
		log.Error(err)
		return "", err
	}
	if err := cd.ping(); err != nil {
		err := fmt.Errorf("cd.ping error {%v}", err)
		log.Error(err)
		return "", err
	}
	token := cd.lockToken(lockid)
	for count := 0; count < cd.info.DistributedLock.MaxRetry; count++ {
		result, err := cd.rdb.SetNX(context.TODO(), lockid, token, time.Duration(cd.info.DistributedLock.TimeOut)*time.Second).Result()
		if err != nil {
			err := fmt.Errorf("redis rdb setnx error {%v}", err)
			log.Error(err)
			return "", err
		}
		if result {
			return token, nil
		}
		rinterval := time.Duration(rand.Intn(cd.info.DistributedLock.RIntervalMs)) * time.Millisecond
		time.Sleep(rinterval)
	}
	err := fmt.Errorf("max number {%v} retry has been exhausted error", cd.info.DistributedLock.MaxRetry)
	log.Error(err)
	return "", err
}

// UnLock 分布式锁解锁
func (cd *Driver) UnLock(lockid string, token string) error {
	var result string
	var err error
	if result, err = cd.rdb.Get(context.TODO(), lockid).Result(); err == nil {
		if result == token {
			if _, err = cd.rdb.Del(context.TODO(), lockid).Result(); err != nil {
				err = fmt.Errorf("redis rdb del error {%v}", err)
				log.Error(err)
				return err
			}
			return nil
		}
		err = fmt.Errorf("redis rdb get lockid {%v} value {%v} != input value {%v} error", lockid, result, token)
		log.Error(err)
		return err
	}
	err = fmt.Errorf("redis rdb get error {%v}", err)
	log.Error(err)
	return err
}
