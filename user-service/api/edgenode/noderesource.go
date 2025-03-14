package edgenode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func listnoderesource(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 获取redis
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprint("edgenodes:resource:*")).Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	edgenodeslist := response.EdgeNodeResourceList{
		EdgeNodes: []response.EdgeNodeResource{},
	}
	for _, tkey := range tkeys {
		clientid := tkey[strings.LastIndex(tkey, ":")+1:]
		enmap, err := rdb.HGetAll(context.TODO(), tkey).Result()
		if err != nil {
			err = fmt.Errorf("redis hgetall error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		for edgenodename, envaluestr := range enmap {
			envalue := &value.EdgeNodeResourceInfo{}
			if err := json.Unmarshal([]byte(envaluestr), envalue); err != nil {
				err = fmt.Errorf("json.Unmarshal error {%v}", err)
				log.Error(err)
				c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
				return
			}

			edgenodeslist.EdgeNodes = append(edgenodeslist.EdgeNodes, response.EdgeNodeResource{
				Name:         edgenodename,
				CpuAll:       envalue.CpuAll,
				MemAll:       envalue.MemAll,
				CpuUse:       envalue.CpuUse,
				MemUse:       envalue.MemUse,
				NvidiaGpuAll: envalue.NvidiaGpuAll,
				ClientID:     clientid,
			})
		}
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: edgenodeslist})
}
