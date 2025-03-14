package user

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"linklab/device-control-v2/login-authentication/api/login"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func userregister(c *gin.Context) {

	var p request.RegisterParameter
	if err := c.ShouldBindJSON(&p); err != nil {
		err = fmt.Errorf("bing json parameter errorcode {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if err := login.UserRegister(p); err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	rmsg := fmt.Sprintf("userid {%v} register success", p.ID)
	log.Info(rmsg)
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: rmsg})
}
