package login

import (
	"context"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/saasbackend"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckUsertTenant(userId string, siteInfo *saasbackend.SiteResult, timeout int) error {

	// 检查租户
	tenantFilter := table.TenantFilterByID{TenantID: siteInfo.ID}
	tenantElem := &table.Tenant{}
	if err := database.Mdriver.FindOneElem("tenants", tenantFilter, tenantElem); err != nil {
		// 若不存在，则创建租户
		tenantElem.TenantID = siteInfo.ID
		tenantElem.TenantName = siteInfo.Name
		tenantElem.IsSystemTenant = false
		if _, err = database.Mdriver.InsertElemIfNotExist("tenants", tenantFilter, tenantElem); err != nil {
			err = fmt.Errorf("database.Mdriver.InsertElemIfNotExist error {%v}", err)
			log.Error(err)
			return err
		}
	} else {
		// 若存在，则检查字段
		if tenantElem.TenantName != siteInfo.Name {
			if _, err := database.Mdriver.UpdateElem("tenants", tenantFilter, bson.M{"tenantName": siteInfo.Name}); err != nil {
				err = fmt.Errorf("database.Mdriver.UpdateElem error {%v}", err)
				log.Error(err)
				return err
			}
		}
	}

	// 检查用户
	userFilter := table.UserFilter{UserID: userId}
	userElem := &table.User{}
	if err := database.Mdriver.FindOneElem("users", userFilter, userElem); err != nil {
		err = fmt.Errorf("database.Mdriver.FindOneElem error {%v}", err)
		log.Error(err)
		return err
	}
	if userElem.TenantID != siteInfo.ID {
		// 更新用户所属租户
		if _, err := database.Mdriver.UpdateElem("users", userFilter, bson.M{"tenantId": siteInfo.ID}); err != nil {
			err = fmt.Errorf("database.Mdriver.UpdateElem error {%v}", err)
			log.Error(err)
			return err
		}
	}

	// 检查缓存
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return err
	}
	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:%s:token:*", "users", userId)).Result()
	if len(tkeys) > 0 {
		if err == nil {
			loginvalue, err := rdb.Get(context.TODO(), tkeys[0]).Result()
			if err != nil {
				err = fmt.Errorf("{%s} redis get {%s} error {%v}", userId, tkeys[0], err)
				log.Error(err)
				return err
			}

			loginstatus := value.UserLoginStatus{}
			if err = json.Unmarshal([]byte(loginvalue), &loginstatus); err != nil {
				err := fmt.Errorf("json.Unmarshal error {%v}", err)
				log.Error(err)
				return err
			}

			if loginstatus.TenantID != siteInfo.ID {
				loginstatus.TenantID = siteInfo.ID
			} else {
				return nil
			}

			loginvaluebyte, err := json.Marshal(loginstatus)
			if err != nil {
				err := fmt.Errorf("json.Marshal error {%v}", err)
				log.Error(err)
				return err
			}

			_, err = rdb.Set(context.TODO(), tkeys[0], loginvaluebyte, time.Second*time.Duration(timeout)).Result()
			if err != nil {
				err = fmt.Errorf("{%v} redis set error {%v}", userId, err)
				log.Error(err)
				return err
			}
		}
	}

	return nil
}
