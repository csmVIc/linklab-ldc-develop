package user

import (
	"errors"
	"fmt"

	"linklab/device-control-v2/base-library/auth"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	uinfo    *UInfo
	ahandler *auth.Handler
)

// RouterInit 用户操作路由初始化
func RouterInit(r *gin.RouterGroup, ui *UInfo) error {
	// 初始化参数传递
	if ui == nil {
		err := errors.New("RouterInit ui is nil")
		log.Error(err)
		return err
	}
	uinfo = ui
	// 登录验证初始化
	var err error
	ahandler, err = auth.New("users")
	if err != nil {
		err = fmt.Errorf("auth New error {%v}", err)
		log.Error(err)
		return err
	}
	// 用户登录，退出登录，用户创建，用户注销
	r.POST("/login", loginhandler)
	r.POST("/signout", ahandler.AuthorizationRequired(), signouthandler)
	r.POST("/register", registerhandler)
	r.POST("/delete", ahandler.AuthorizationRequired(), deletehandler)
	return nil
}

func loginhandler(c *gin.Context) {
	userlogin(c, uinfo.Login.TimeOut, uinfo.Register.Email)
}

func signouthandler(c *gin.Context) {
}

func registerhandler(c *gin.Context) {
	userregister(c)
}

func deletehandler(c *gin.Context) {
}
