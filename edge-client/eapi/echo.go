package eapi

import (
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func echo(c *gin.Context) {
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
}
