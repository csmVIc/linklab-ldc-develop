package auth

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"linklab/device-control-v2/base-library/parameter/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	// Adriver 全局验证实例
	Adriver *Handler
)

// Handler 用户登录验证句柄
type Handler struct {
	// uType 用户类型 users, clients等.
	uType string
}

// New 产生新的实例
func New(uType string) (*Handler, error) {
	if len(uType) < 1 {
		err := errors.New("New len(uType) < 1 error")
		log.Error(err)
		return nil, err
	}
	ah := &Handler{uType: uType}
	return ah, nil
}

// AuthorizationRequired 用户登录验证
func (ah *Handler) AuthorizationRequired() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 先检查是否为websocket连接
		token := ah.getWebsocketToken(c)
		if len(token) < 1 {
			// 如果不是websocket连接
			token = c.Request.Header.Get("Authorization")
			if len(token) < 1 {
				err := fmt.Errorf("interface {%v} require authorization, field {Authorization} need in header", c.Request.URL.RequestURI())
				log.Errorln(err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: -1, Msg: err.Error()})
				return
			}
		}

		issha256, err := regexp.MatchString("^[A-Fa-f0-9]{64}$", token)
		if err != nil || issha256 == false {
			err = fmt.Errorf("token {%v} not be sha256 error", token)
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		id, err := ah.checkLoginStatus(token)
		if err != nil {
			err = fmt.Errorf("{%v} login status of token {%v} not exist error {%v}", ah.uType, token, err)
			log.Errorln(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		log.Infof("{%v} pass {%v} authentication", id, ah.uType)
		c.Set("id", id)
		c.Next()
	}
}
