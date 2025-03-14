package sub

import (
	"context"
	"crypto/sha256"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func (sd *Driver) createGroupID(userid string) (string, error) {

	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return "", err
	}

	for count := 0; count < sd.info.MaxCreateGroupIDRetry; count++ {
		ntime := time.Now().Unix()
		rint := rand.Int()
		groupidseed := fmt.Sprintf("%s:%v:%v", userid, ntime, rint)
		groupidbinary := sha256.Sum256([]byte(groupidseed))
		groupid := fmt.Sprintf("%x", groupidbinary)

		gexist, err := rdb.HExists(context.TODO(), "pods:apply", groupid).Result()
		if err != nil {
			err = fmt.Errorf("rdb.HExists error {%v}", err)
			log.Error(err)
			return "", err
		}

		if gexist == false {
			return groupid, nil
		}
	}

	err = fmt.Errorf("max number {%v} of create groupid has been exhausted", sd.info.MaxCreateGroupIDRetry)
	log.Error(err)
	return "", err
}
