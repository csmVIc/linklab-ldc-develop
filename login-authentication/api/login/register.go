package login

import (
	"crypto/sha256"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/request"

	log "github.com/sirupsen/logrus"
)

// UserRegister 用户注册
func UserRegister(p request.RegisterParameter) error {

	// 随机生成用户盐值
	salt := auth.CreateSalt(p.ID)

	// 计算最终存储的哈希值
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v%v", p.Password, salt))))

	if p.TenantID < 0 {
		// 如果指定的tenantid为负值,则需要查找系统默认的tenantid
		filter := &table.TenantFilterByFlag{
			IsSystemTenant: true,
		}
		tenant := &table.Tenant{}
		if err := database.Mdriver.FindOneElem("tenants", filter, tenant); err != nil {
			err = fmt.Errorf("database find default tenant error {%v}", err)
			log.Error(err)
			return err
		}
		p.TenantID = tenant.TenantID
	} else {
		// 查找p.TenantID是否存在
		filter := &table.TenantFilterByID{
			TenantID: p.TenantID,
		}
		if err := database.Mdriver.DocExist("tenants", filter); err != nil {
			err = fmt.Errorf("database not find tenant {%v} error {%v}", p.TenantID, err)
			log.Error(err)
			return err
		}
	}

	// 用户插入到数据库中
	user := &table.User{
		UserID:   p.ID,
		Email:    p.Email,
		Hash:     hash,
		Salt:     salt,
		TenantID: p.TenantID,
	}
	filter := &table.UserFilter{
		UserID: p.ID,
	}

	// 如果用户名或者密码未被占用,则注册该用户
	result, err := database.Mdriver.InsertElemIfNotExist("users", filter, user)
	if err != nil {
		err = fmt.Errorf("database.Mdriver.InsertElemIfNotExist error {%v}", err)
		return err
	}

	// 检查是否重复注册
	if result.UpsertedCount < 1 {
		err = fmt.Errorf("userid {%v} or email {%v} already registered", p.ID, p.Email)
		log.Error(err)
		// 考虑到webide存在重复注册的可能性,取消返回错误
		// c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		// return
	}

	return nil
}

// JudgeUserRegisterd 判断用户是否注册过
func JudgeUserRegisterd(userid string) error {

	filter := &table.UserFilter{
		UserID: userid,
	}

	return database.Mdriver.DocExist("users", filter)
}
