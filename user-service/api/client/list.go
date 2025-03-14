package client

import (
	"context"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func listclient(c *gin.Context) {

	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// p := request.ClientListQuery{}
	// if err := c.ShouldBindQuery(&p); err != nil {
	// 	err = fmt.Errorf("bing query parameter error {%v}", err)
	// 	log.Error(err)
	// 	c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
	// 	return
	// }

	// 获取redis
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:*:token:*", "clients")).Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	clients := response.ClientList{
		Clients: []response.ClientStatus{},
	}
	for _, tkey := range tkeys {
		clientid, err := auth.GetIDFromKey(tkey)
		if err != nil {
			err = fmt.Errorf("auth.GetIDFromKey error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		tenantid, err := cache.Cdriver.GetClientTenantID(clientid)
		if err != nil {
			err = fmt.Errorf("cache.Cdriver.GetClientTenantID error {%v}", err)
			log.Error(err)
			// 可能某个客户端退出，因此不返回错误，直接跳过
			continue
		}

		tenantids := []int{}
		for key := range tenantid {
			tenantids = append(tenantids, key)
		}

		clients.Clients = append(clients.Clients, response.ClientStatus{
			ClientID: clientid,
			TenantID: tenantids,
		})
	}
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: clients})
}
