package fcache

import (
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func imagebuilddownload(c *gin.Context) {
	// 获取用户id
	id := c.GetString("id")
	if len(id) < 1 {
		err := errors.New("id not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// http 参数验证
	var p request.ImageBuildDownload
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("{%v} bind query parameter error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 查找文件
	filter := &table.ImageBuildFilter{
		FileHash: p.FileHash,
	}
	result := &table.ImageBuild{}
	log.Debugf("imagebuild download filter {%v}", filter)
	if err := database.Mdriver.FindOneElem("imagebuild", filter, result); err != nil {
		// 需要的文件不存在
		err = fmt.Errorf("{%v} imagebuild filehash {%v} not exist error {%v}", id, p.FileHash, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	log.Infof("{%v} download imagebuild {%v} bytes", id, len(result.FileData.Data))

	// 需要的文件存在
	c.Data(http.StatusOK, "application/octet-stream", result.FileData.Data)
}
