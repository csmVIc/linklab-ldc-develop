package fcache

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func imagebuildupload(c *gin.Context) {
	// 获取用户id
	id := c.GetString("id")
	if len(id) < 1 {
		err := errors.New("id not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// 文件读取
	fileheader, err := c.FormFile("file")
	if err != nil {
		err = fmt.Errorf("{%v} read upload form file header error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	f, err := fileheader.Open()
	if err != nil {
		err = fmt.Errorf("{%v} open upload file error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	binary, err := ioutil.ReadAll(f)
	if err != nil {
		err = fmt.Errorf("{%v} read upload file binary error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// 计算文件hash值
	filehash := fmt.Sprintf("%x", sha256.Sum256(binary))
	// 若数据库中不存在该文件，则插入文件
	elem := &table.ImageBuild{
		FileHash: filehash,
		FileData: primitive.Binary{
			Subtype: 0x00,
			Data:    binary,
		},
	}
	filter := &table.ImageBuildFilter{
		FileHash: filehash,
	}
	_, err = database.Mdriver.InsertElemIfNotExist("imagebuild", filter, elem)
	if err != nil {
		err = fmt.Errorf("{%v} database insert error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 返回结果
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: map[string]string{
		"filehash": filehash,
	}})
}
