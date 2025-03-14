package login

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"linklab/device-control-v2/base-library/auth"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/cache/value"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/parameter/request"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

// Authentication 身份验证
func Authentication(p *request.LoginParameter, utype string, timeout int) (string, error) {

	// 日志记录
	ispass := false
	defer func() {
		tags := map[string]string{
			"userid": p.ID,
		}
		fields := map[string]interface{}{
			"ispass":    ispass,
			"logintime": time.Now().UnixNano(),
			"podname":   os.Getenv("POD_NAME"),
			"nodename":  os.Getenv("NODE_NAME"),
		}
		err := logger.Ldriver.WriteLog("userlogin", tags, fields)
		if err != nil {
			err = fmt.Errorf("write log err {%v}", err)
			log.Error(err)
		}
	}()

	// 密码需要符合sha256的regex
	issha256, err := regexp.MatchString("^[A-Fa-f0-9]{64}$", p.Password)
	if err != nil || issha256 == false {
		err = fmt.Errorf("password not be sha256 error")
		log.Error(err)
		return "", err
	}

	// 获取缓存
	rdb, err := cache.Cdriver.GetRdb()
	if err != nil {
		err = fmt.Errorf("redis get rdb error {%s}", err)
		log.Error(err)
		return "", err
	}

	// 检查是否已经登录
	tkeys, err := rdb.Keys(context.TODO(), fmt.Sprintf("%s:id:%s:token:*", utype, p.ID)).Result()
	if err == nil {
		if len(tkeys) == 1 {

			value, err := rdb.Get(context.TODO(), tkeys[0]).Result()
			if err != nil {
				log.Errorf("{%s} redis get {%s} error {%v}", p.ID, tkeys[0], err)
				return "", err
			}

			trueCheck, salt, err := getTrueCheckAndSaltFromValue(value)
			if err != nil {
				log.Errorf("getTrueCheckAndSaltFromValue error {%v}", err)
				return "", err
			}

			computeCheck := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v%v", p.Password, salt))))
			if computeCheck != trueCheck {
				// 已经登录但是本次登录,密码输入错误
				log.Errorf("{%v} compute check {%v} != true check {%v}", p.ID, computeCheck, trueCheck)
				return "", fmt.Errorf("{%v} password check error", p.ID)
			}

			// 已经登录则只需要延长token的时间
			log.Infof("%s {%s} already logged in, token {%v}\n", utype, p.ID, tkeys[0])
			_, err = rdb.Expire(context.TODO(), tkeys[0], time.Second*time.Duration(timeout)).Result()
			if err != nil {
				log.Errorf("{%s} redis expire error {%v}", p.ID, err)
				return "", err
			}

			// 解析出token
			token, err := auth.GetTokenFromKey(tkeys[0])
			if err != nil {
				log.Errorf("{%s} get token from error {%v}", tkeys[0], err)
				return "", err
			}
			ispass = true
			return token, nil
		} else if len(tkeys) > 1 {
			// 如果发现多个token,则系统可能出现问题
			log.Errorf("%s {%s} already logged in, has multiple token error {%v}\n", utype, p.ID, tkeys)
			return "", fmt.Errorf("%s {%s} already logged in, has multiple token error", utype, p.ID)
		}
	}

	// 未登录过
	trueCheck, salt, tenantid, err := getTrueCheckAndSaltAndTenantIDFromDb(p.ID, utype)
	if err != nil {
		err := fmt.Errorf("{%v} {%v} get password and salt from database error {%v}", utype, p.ID, err)
		log.Error(err)
		return "", err
	}

	// 计算校验值的计算值
	computeCheck := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%v%v", p.Password, salt))))
	// 校验值的真值和计算值相等
	if computeCheck == trueCheck {

		// 计算token值
		token := auth.CreateToken(p.ID, computeCheck)

		// valuestr
		loginstatus := value.UserLoginStatus{
			TrueCheck: trueCheck,
			Salt:      salt,
			TenantID:  tenantid,
		}
		valuebyte, err := json.Marshal(loginstatus)
		if err != nil {
			err := fmt.Errorf("json.Marshal error {%v}", err)
			log.Error(err)
			return "", err
		}

		_, err = rdb.Set(context.TODO(), fmt.Sprintf("%s:id:%s:token:%s", utype, p.ID, token), string(valuebyte), time.Second*time.Duration(timeout)).Result()
		if err != nil {
			err = fmt.Errorf("{%v} redis set error {%v}", p.ID, err)
			log.Error(err)
			return "", err
		}

		log.Infof("{%v} login success", p.ID)
		ispass = true
		return token, nil
	}
	// 校验值的真值和计算值不相等
	log.Error("{%v} compute check {%v} != true check {%v}", p.ID, computeCheck, trueCheck)
	return "", fmt.Errorf("{%v} password check error", p.ID)
}
