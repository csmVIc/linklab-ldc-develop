package fcache

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func fileupload(c *gin.Context) {
	// 获取用户id
	id := c.GetString("id")
	if len(id) < 1 {
		err := errors.New("id not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	// HTTP POST参数验证
	pStr := c.PostForm("parameters")
	var p request.FileUpload
	if err := json.Unmarshal([]byte(pStr), &p); err != nil || len(p.BoardName) < 1 {
		err = fmt.Errorf("{%v} bing json parameter error {%v}", id, err)
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
	elem := &table.File{
		BoardName: p.BoardName,
		FileHash:  filehash,
		FileData: primitive.Binary{
			Subtype: 0x00,
			Data:    binary,
		},
	}
	filter := &table.FileFilter{
		BoardName: p.BoardName,
		FileHash:  filehash,
	}
	log.Infof("fileupload info filehash{%v} boardname{%v}", elem.FileHash, elem.BoardName)
	// 若数据库中不存在该文件,则插入文件
	result, err := database.Mdriver.InsertElemIfNotExist("files", filter, elem)
	if err != nil {
		err = fmt.Errorf("{%v} database insert error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	isfind := false
	if result.UpsertedCount < 1 {
		isfind = true
	}

	// 日志记录
	tags := map[string]string{
		"userid": id,
	}
	fields := map[string]interface{}{
		"boardname":  p.BoardName,
		"filehash":   filehash,
		"isfind":     isfind,
		"bytes":      len(binary),
		"uploadtime": time.Now().UnixNano(),
		"podname":    os.Getenv("POD_NAME"),
		"nodename":   os.Getenv("NODE_NAME"),
	}
	err = logger.Ldriver.WriteLog("fileupload", tags, fields)
	if err != nil {
		err = fmt.Errorf("write log err {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 返回结果
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: map[string]string{
		"filehash": filehash,
	}})
}
