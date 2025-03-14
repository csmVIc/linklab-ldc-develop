package auth

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// EdgeClientAuthorizationRequired 验证
func EdgeClientAuthorizationRequired(cloudtoken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if len(token) < 1 {
			err := fmt.Errorf("interface {%v} require authorization, field {Authorization} need in header", c.Request.URL.RequestURI())
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		issha256, err := regexp.MatchString("^[A-Fa-f0-9]{64}$", token)
		if err != nil || issha256 == false {
			err = fmt.Errorf("token {%v} not be sha256 error", token)
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		if token != cloudtoken {
			err = fmt.Errorf("token {%v} != cloud send token error", token)
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		c.Next()
	}
}
