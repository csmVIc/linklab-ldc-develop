package fcache

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/logger"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func filedownload(c *gin.Context) {
	// 获取用户id
	id := c.GetString("id")
	if len(id) < 1 {
		err := errors.New("id not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// http 参数验证
	var p request.FileDownload
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("{%v} bind query parameter error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// 查找文件
	filter := &table.FileFilter{
		BoardName: p.BoardName,
		FileHash:  p.FileHash,
	}
	result := &table.File{}
	log.Debugf("filedownload file filter {%v}", filter)
	if err := database.Mdriver.FindOneElem("files", filter, result); err != nil {
		// 需要的文件不存在
		err = fmt.Errorf("{%v} board {%v} filehash {%v} not exist error {%v}", id, p.BoardName, p.FileHash, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	log.Infof("{%v} download {%v} file {%v} bytes", id, result.BoardName, len(result.FileData.Data))

	// 日志记录
	tags := map[string]string{
		"clientid": id,
	}
	fields := map[string]interface{}{
		"boardname":    p.BoardName,
		"filehash":     p.FileHash,
		"bytes":        len(result.FileData.Data),
		"downloadtime": time.Now().UnixNano(),
		"podname":      os.Getenv("POD_NAME"),
		"nodename":     os.Getenv("NODE_NAME"),
	}
	err := logger.Ldriver.WriteLog("filedownload", tags, fields)
	if err != nil {
		err = fmt.Errorf("write log err {%v}", err)
		log.Error(err)
		return
	}

	// 需要的文件存在
	c.Data(http.StatusOK, "application/octet-stream", result.FileData.Data)
}
