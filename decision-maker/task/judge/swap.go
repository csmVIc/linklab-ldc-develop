package judge

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func (jd *Driver) taskQueueSwap(rdb *redis.ClusterClient) error {

	// 触发交换
	tasksQueueLen, err := rdb.LLen(context.TODO(), "tasks:queue").Result()
	if err != nil {
		err = fmt.Errorf("redis llen error {%s}", err)
		log.Error(err)
		return err
	}
	if tasksQueueLen < 1 {
		for {
			tasksTmpLen, err := rdb.LLen(context.TODO(), "tasks:tmp").Result()
			if err != nil {
				err = fmt.Errorf("redis llen error {%s}", err)
				log.Error(err)
				return err
			}
			if tasksTmpLen < 1 {
				break
			}
			if _, err := rdb.BRPopLPush(context.TODO(), "tasks:tmp", "tasks:queue", time.Second).Result(); err != nil {
				err = fmt.Errorf("redis BRPopLPush error {%s}", err)
				log.Error(err)
				return err
			}
		}
	}

	return nil
}
