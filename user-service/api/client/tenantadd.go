package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/saasbackend"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func addclienttenant(c *gin.Context) {
	// 获取用户id
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP POST 参数验证
	p := request.ClientTenantSet{}
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查tenantid是否存在
	tidfilter := table.TenantFilterByID{
		TenantID: p.TenantID,
	}
	if err := database.Mdriver.DocExist("tenants", tidfilter); err != nil {
		// 检查saas的接口
		siteInfo, err := saasbackend.SDriver.GetSiteInfo(p.TenantID)
		if err != nil {
			// 返回tenantid不存在
			err = fmt.Errorf("tenantid not exist error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
		// 创建租户
		tenantFilter := table.TenantFilterByID{TenantID: siteInfo.ID}
		tenantElem := &table.Tenant{
			TenantID:       siteInfo.ID,
			TenantName:     siteInfo.Name,
			IsSystemTenant: false,
		}
		if _, err = database.Mdriver.InsertElemIfNotExist("tenants", tenantFilter, tenantElem); err != nil {
			err = fmt.Errorf("database.Mdriver.InsertElemIfNotExist error error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 获取redis
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查设备管理客户端是否活跃
	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:%s:token:*", "clients", p.ClientID)).Result()
	if err != nil {
		err = fmt.Errorf("redis keys error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	if len(tkeys) < 1 {
		err = fmt.Errorf("clientid {%v} not active error {%v}", p.ClientID, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 检查tenantid是否已经存在
	ptid, err := cache.Cdriver.GetClientTenantID(p.ClientID)
	if err != nil {
		err = fmt.Errorf("get pre client tenantid error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	if _, isOk := ptid[p.TenantID]; isOk {
		c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
		return
	}

	// 增加redis中的tenantid
	loginstatus := value.ClientLoginStatus{}
	if err := cache.Cdriver.GetLoginStatus("clients", p.ClientID, &loginstatus); err != nil {
		err = fmt.Errorf("cache.Cdriver.GetLoginStatus error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	loginstatus.TenantID[p.TenantID] = true
	lsvalue, err := json.Marshal(loginstatus)
	if err != nil {
		err = fmt.Errorf("json.Marshal error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	_, err = rdb.Set(context.TODO(), tkeys[0], string(lsvalue), time.Second*time.Duration(cinfo.ClientTenant.ClientCacheTTL)).Result()
	if err != nil {
		err = fmt.Errorf("rdb.Set error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 增加mongodb中的tenantid
	clientfilter := table.ClientFilter{
		UserName: p.ClientID,
	}
	ptid[p.TenantID] = true
	if _, err = database.Mdriver.UpdateElem("clients", clientfilter, bson.M{
		"tenantId": ptid,
	}); err != nil {
		err = fmt.Errorf("database.Mdriver.UpdateElem error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
	return
}
