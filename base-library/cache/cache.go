package cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

// Driver cache客户端
type Driver struct {
	info *CInfo
	rdb  *redis.ClusterClient
}

var (
	// Cdriver cache客户端全局实例
	Cdriver *Driver
)

func (cd *Driver) ping() error {
	if cd.rdb == nil {
		err := errors.New("cd.rdb is nil error")
		log.Error(err)
		return err
	}
	pong, err := cd.rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Errorf("ping error %v %v", pong, err)
	}
	return err
}

func (cd *Driver) init() error {

	clusterSlots := func(ctx context.Context) ([]redis.ClusterSlot, error) {
		nodes := []redis.ClusterNode{}
		for _, addr := range cd.info.Client.Address {
			nodes = append(nodes, redis.ClusterNode{
				Addr: fmt.Sprintf("%s:%s", addr.Host, addr.Port),
			})
		}
		slots := []redis.ClusterSlot{
			{
				Start: 0,
				End:   16383,
				Nodes: nodes,
			},
		}
		return slots, nil
	}

	cd.rdb = redis.NewClusterClient(&redis.ClusterOptions{
		ClusterSlots:  clusterSlots,
		RouteRandomly: true,
		Password:      cd.info.Client.PassWord,
	})

	if err := cd.ping(); err != nil {
		log.Errorf("cd.ping() error {%v}", err)
		return err
	}

	return nil
}

// GetRdb 返回redis数据库实例
func (cd *Driver) GetRdb() (*redis.ClusterClient, error) {
	if err := cd.ping(); err != nil {
		if err = cd.init(); err != nil {
			log.Errorf("cd.init() error %s", err)
			return nil, err
		}
	}
	return cd.rdb, nil
}

// New 创建执行实例
func New(i *CInfo) (*Driver, error) {
	if i == nil {
		return nil, errors.New("init info nil error")
	}
	id := &Driver{info: i, rdb: nil}
	id.init()
	if err := id.ping(); err != nil {
		log.Errorf("id.ping() error %v", err)
		return nil, err
	}
	return id, nil
}
