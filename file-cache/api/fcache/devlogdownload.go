package fcache

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func devlogdownload(c *gin.Context) {
	// 获取用户id
	id := c.GetString("id")
	if len(id) < 1 {
		err := errors.New("id not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// HTTP参数验证
	p := request.DevLogDownload{}
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("{%v} bind query parameter error {%v}", id, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 过滤器
	filter := bson.M{
		"userId":  id,
		"groupId": p.GroupID,
	}
	if len(p.ClientID) > 0 && len(p.DeviceID) > 0 {
		filter["clientId"] = p.ClientID
		filter["devPort"] = p.DeviceID
	}
	log.Debugf("devlogdownload devicelog filter {%v}", filter)

	// 查找日志
	cursor, err := database.Mdriver.FindElem("devicelog", filter)
	if err != nil {
		// 需要的日志不存在
		err = fmt.Errorf("devlog {%v} not exist error {%v}", filter, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	arr := []table.DeviceLog{}
	if err := cursor.All(context.TODO(), &arr); err != nil {
		err = fmt.Errorf("cursor.All error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if len(arr) < 1 {
		err = fmt.Errorf("find devlog length {%v} is 0 error", len(arr))
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	zipbuf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(zipbuf)
	for index, devlog := range arr {

		// 创建压缩文件
		fname := fmt.Sprintf("device-%v.json", index)
		zipFile, err := zipWriter.Create(fname)
		if err != nil {
			err = fmt.Errorf("zipWriter.Create {%v} error {%v}", fname, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 生成文件内容
		body := response.DeviceLog{
			ClientID: devlog.ClientID,
			DeviceID: devlog.DevPort,
			Logs:     devlog.Logs,
		}
		bodybytes, err := json.MarshalIndent(body, "", "		")
		if err != nil {
			err = fmt.Errorf("json.Marshal error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}

		// 写入文件内容
		_, err = zipFile.Write(bodybytes)
		if err != nil {
			err = fmt.Errorf("zipFile.Write error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 关闭文件
	err = zipWriter.Close()
	if err != nil {
		err = fmt.Errorf("zipWriter.Close error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", zipbuf.Bytes())
}
