package edgenode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/tool"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func listpod(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 命名空间
	namespace := tool.CreateMD5(userid)

	// 参数验证
	p := request.UserPodQuery{}
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
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

	tkeys, err := rdb.Keys(context.TODO(), "pods:active:*").Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	podslist := response.PodList{
		Pods: []response.PodInfo{},
	}
	for _, tkey := range tkeys {
		clientid := tkey[strings.LastIndex(tkey, ":")+1:]
		podmap, err := rdb.HGetAll(context.TODO(), tkey).Result()
		if err != nil {
			err = fmt.Errorf("redis hgetall error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		for _, podvaluestr := range podmap {

			podvalue := &value.PodInfo{}
			if err := json.Unmarshal([]byte(podvaluestr), podvalue); err != nil {
				err = fmt.Errorf("json.Unmarshal error {%v}", err)
				log.Error(err)
				c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
				return
			}

			if p.AllPods || (p.AllPods == false && podvalue.Namespace == namespace) {
				podelem := response.PodInfo{
					Name:       podvalue.Name,
					Namespace:  podvalue.Namespace,
					NodeName:   podvalue.NodeName,
					Ready:      podvalue.Ready,
					ClientID:   clientid,
					CreateTime: podvalue.CreateTime,
					Containers: []response.ContainerInfo{},
				}
				for _, containervalue := range podvalue.Containers {
					containerelem := response.ContainerInfo{
						Name:         containervalue.Name,
						Ready:        containervalue.Ready,
						RestartCount: containervalue.RestartCount,
						Image:        containervalue.Image,
					}
					podelem.Containers = append(podelem.Containers, containerelem)
				}
				podslist.Pods = append(podslist.Pods, podelem)
			}
		}
	}
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: podslist})
}
