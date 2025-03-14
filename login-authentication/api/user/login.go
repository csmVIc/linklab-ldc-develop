package user

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/base-library/saasbackend"
	"linklab/device-control-v2/login-authentication/api/login"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func userlogin(c *gin.Context, timeout int, emailpattern string) {

	// http post 参数验证
	var p request.LoginParameter
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 首先判断用户是否注册过
	if err := login.JudgeUserRegisterd(p.ID); err != nil {
		// 若未注册，那么默认注册
		if err := login.UserRegister(request.RegisterParameter{
			ID:       p.ID,
			Password: p.Password,
			Email:    fmt.Sprintf(emailpattern, p.ID),
			TenantID: 1,
		}); err != nil {
			err = fmt.Errorf("user {%v} register error {%v}", p.ID, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 获取用户的租户信息
	siteInfo, err := saasbackend.SDriver.GetUserSiteInfo(p.ID)
	if err != nil {
		log.Errorf("user {%v} GetUserSiteInfo error {%v}", p.ID, err)
	}

	// 检查用户的租户信息
	if siteInfo != nil {
		if err = login.CheckUsertTenant(p.ID, siteInfo, timeout); err != nil {
			err = fmt.Errorf("user {%v} CheckUsertTenant error {%v}", p.ID, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 登录验证
	token, err := login.Authentication(&p, "users", timeout)
	if err != nil {
		err = fmt.Errorf("user {%v} login auth error {%v}", p.ID, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 返回token
	msg := fmt.Sprintf("user {%v} login success", p.ID)
	log.Info(msg)
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: msg, Data: map[string]string{
		"token": token,
	}})
	return
}
