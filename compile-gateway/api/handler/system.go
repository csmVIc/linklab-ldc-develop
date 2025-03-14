package handler

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/messenger"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CompileSystem 编译系统接口
func CompileSystem(c *gin.Context) {
	pStr := c.PostForm("parameters")
	p := request.CompileSystemUpload{}
	if err := json.Unmarshal([]byte(pStr), &p); err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: "bind json request error"})
		return
	}

	fileheader, err := c.FormFile("patch")
	if err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: "parse patch file error"})
		return
	}

	f, err := fileheader.Open()
	if err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: "open patch file error"})
		return
	}

	binary, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: "read patch file error"})
		return
	}

	// 哈希值计算
	hashhandler := sha1.New()
	len, err := hashhandler.Write(binary)
	if len < 1 || err != nil {
		err := fmt.Errorf("hash sha1 write error {%v} || binary write length < 1 {%v}", len, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	sha1hex := hex.EncodeToString(hashhandler.Sum(nil))
	if sha1hex != p.FileHash {
		err := fmt.Errorf("hash sha1 compute {%v} != input {%v}", sha1hex, p.FileHash)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	elem := &table.CompileTable{
		Type:        "system",
		CompileType: p.CompileType,
		BoardType:   p.BoardType,
		FileHash:    p.FileHash,
		Branch:      p.Branch,
		FileData: primitive.Binary{
			Subtype: 0x00,
			Data:    binary,
		},
		Output: primitive.Binary{
			Subtype: 0x00,
			Data:    []byte{},
		},
		Message: "",
		Status:  "input",
	}
	filter := &table.CompileTableFilter{
		CompileType: p.CompileType,
		BoardType:   p.BoardType,
		FileHash:    p.FileHash,
	}
	updateresult, err := database.Mdriver.InsertElemIfNotExist("compile", filter, elem)
	if err != nil {
		err := fmt.Errorf("driver.Md.InsertElemIfNotExist error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if updateresult.UpsertedCount > 0 {
		// 之前不存在，现在插入成功
		reply := msg.ReplyMsg{}
		task := msg.CompileTask{
			CompileType: p.CompileType,
			BoardType:   p.BoardType,
			FileHash:    p.FileHash,
		}
		if err := messenger.Mdriver.RequestMsg(hinfo.CompileSystem.Topic+"."+p.CompileType, task, time.Duration(hinfo.CompileSystem.ReplyTimeout)*time.Second, &reply); err != nil {
			// request出错
			err = fmt.Errorf("RequestMsg error {%v}", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
		if reply.Code < 0 {
			err = fmt.Errorf("RequestMsg reply error {%v}", reply.Msg)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
		c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success"})
		return
	} else if updateresult.MatchedCount > 0 {
		// 之前已存在
		c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "task already exists"})
		return
	}

	err = fmt.Errorf("update result status {%v} not exist error", updateresult)
	c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
	return
}
