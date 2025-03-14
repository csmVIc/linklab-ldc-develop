package handler

import (
	"errors"
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

// DownloadNonBlockhandler 非阻塞下载编译结果函数
func DownloadNonBlockhandler(c *gin.Context) {

	var p request.CompileDownload
	if err := c.ShouldBindQuery(&p); err != nil {
		log.Error("/compile/nonblock bind query error")
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: "/compile/nonblock bind query error"})
		return
	}

	filter := table.CompileTableFilter{
		CompileType: p.CompileType,
		BoardType:   p.BoardType,
		FileHash:    p.FileHash,
	}
	result := table.CompileTable{}
	if err := database.Mdriver.FindOneElem("compile", filter, &result); err == nil {
		if result.Status == "output" {
			c.Data(http.StatusOK, "application/octet-stream", result.Output.Data)
			return
		} else if result.Status == "builderr" {
			c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "error", Data: map[string]string{
				"message": result.Message,
			}})
			return
		} else if result.Status == "input" {
			c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "compiling"})
			return
		}
		err := fmt.Errorf("/compile/nonblock Status {%v} not exist error", result.Status)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	} else {
		c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "notexist"})
		return
	}
}

// DownloadBlockhandler 阻塞下载编译结果函数
func DownloadBlockhandler(c *gin.Context) {

	var p request.CompileDownload
	if err := c.ShouldBindQuery(&p); err != nil {
		log.Error("/compile/block bind query error")
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: "/compile/block bind query error"})
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
				c.Data(http.StatusOK, "application/octet-stream", result.Output.Data)
				return
			} else if result.Status == "builderr" {
				c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "error", Data: map[string]string{
					"message": result.Message,
				}})
				return
			}
		} else {
			c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "notexist"})
			return
		}
		time.Sleep(time.Second)
	}

	err := errors.New("/compile/block query timeout")
	log.Error(err)
	c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: fmt.Sprintf("/compile/block download binary timeout %vs error", hinfo.Download.Timeout)})
	return
}
