package handler

import (
	"fmt"

	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetBlockStatusHandler 阻塞获取编译状态
func GetBlockStatusHandler(c *gin.Context) {

	var p request.CompileDownload
	if err := c.ShouldBindQuery(&p); err != nil {
		log.Errorf("/compile/block/status bind query error {%v}", err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: fmt.Sprintf("/compile/block/status bind query error {%v}", err)})
		return
	}

	filter := table.CompileTableFilter{
		CompileType: p.CompileType,
		BoardType:   p.BoardType,
		FileHash:    p.FileHash,
	}
	result := table.CompileTable{}
	for i := int64(0); i < hinfo.Download.Timeout; i++ {
		if err := database.Mdriver.FindOneElem("compile", filter, &result); err == nil {
			if result.Status == "output" {
				c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "completed"})
				return
			} else if result.Status == "builderr" {
				c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "error", Data: map[string]string{
					"message": result.Message,
				}})
				return
			} else if result.Status == "input" {
				log.Debugf("{%v} compiling", p)
			}
		} else {
			c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "notexist"})
			return
		}
		time.Sleep(time.Second)
	}

	err := fmt.Errorf("/compile/block/status query timeout %vs error", hinfo.Download.Timeout)
	log.Error(err)
	c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
	return

}
