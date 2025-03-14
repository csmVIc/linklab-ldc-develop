package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// CreateDeviceBindGroupID 创建设备绑定组ID
func (cd *Driver) CreateDeviceBindGroupID(boardgroup string) (string, error) {

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

	// 查询现有ID计数
	if isOk, err := cd.rdb.HExists(context.TODO(), "bind:group:id", boardgroup).Result(); err != nil {

		err := fmt.Errorf("cd.rdb.HExists error {%v}", err)
		log.Error(err)
		return "", err

	} else if isOk {

		countStr, err := cd.rdb.HGet(context.TODO(), "bind:group:id", boardgroup).Result()
		if err != nil {
			err := fmt.Errorf("cd.rdb.HGet error {%v}", err)
			log.Error(err)
			return "", err
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			err := fmt.Errorf("strconv.Atoi error {%v}", err)
			log.Error(err)
			return "", err
		}

		count += 1

		if _, err := cd.rdb.HSet(context.TODO(), "bind:group:id", boardgroup, count).Result(); err != nil {
			err := fmt.Errorf("cd.rdb.HSet error {%v}", err)
			log.Error(err)
			return "", err
		}

		return fmt.Sprintf("%v-%v", boardgroup, count), nil

	} else {

		if _, err := cd.rdb.HSet(context.TODO(), "bind:group:id", boardgroup, 0).Result(); err != nil {
			err := fmt.Errorf("cd.rdb.HSet error {%v}", err)
			log.Error(err)
			return "", err
		}

		return fmt.Sprintf("%v-%v", boardgroup, 0), nil

	}
}
