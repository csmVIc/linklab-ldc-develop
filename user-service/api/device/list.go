package device

import (
	"context"
	"errors"
	"fmt"
	"linklab/device-control-v2/base-library/cache"
	"linklab/device-control-v2/base-library/database"
	"linklab/device-control-v2/base-library/database/table"
	"linklab/device-control-v2/base-library/parameter/request"
	"linklab/device-control-v2/base-library/parameter/response"
	"net/http"
	"io"
	"linklab/device-control-v2/base-library/tool"
	"github.com/gin-gonic/gin"
	"github.com/docker/distribution/registry/client"
	"strings"
	log "github.com/sirupsen/logrus"
)

func listalldevice(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	var p request.DeviceListQuery
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	if p.BoardName != "all" {
		filter := table.BoardFilter{BoardName: p.BoardName}
		if err := database.Mdriver.DocExist("boards", filter); err != nil {
			err = fmt.Errorf("boardname {%v} not support error {%v}", p.BoardName, err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}

	// 设备加锁
	devlocktoken, err := cache.Cdriver.Lock("devices:lock")
	if err != nil {
		err = fmt.Errorf("cache.Cdriver.Lock {%v} error {%v}", "devices:lock", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	defer func() {
		err = cache.Cdriver.UnLock("devices:lock", devlocktoken)
		if err != nil {
			err = fmt.Errorf("cache.Cdriver.UnLock {%v} error {%v}", "devices:lock", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}()

	devices, err := alldevcies(p.BoardName)
	if err != nil {
		err = fmt.Errorf("get boardname {%v} devices error {%v}", p.BoardName, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: response.DeviceList{Devices: *devices}})
}

func listuserdevice(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 设备加锁
	devlocktoken, err := cache.Cdriver.Lock("devices:lock")
	if err != nil {
		err = fmt.Errorf("cache.Cdriver.Lock {%v} error {%v}", "devices:lock", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}
	defer func() {
		err = cache.Cdriver.UnLock("devices:lock", devlocktoken)
		if err != nil {
			err = fmt.Errorf("cache.Cdriver.UnLock {%v} error {%v}", "devices:lock", err)
			log.Error(err)
			c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
			return
		}
	}()

	devices, err := userdevices(userid)
	if err != nil {
		err = fmt.Errorf("get user {%v} devices error {%v}", userid, err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: *devices})
}

func getcurrtasknum(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	var p request.WaitTime
	if err := c.ShouldBindQuery(&p); err != nil {
		err = fmt.Errorf("bing query parameter error {%v}", err)
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	rdb, _ := cache.Cdriver.GetRdb()

	var reNum int64
	var waitTime int64
	if p.Queuetype == "queue" {
		tasknum, _ := rdb.LLen(context.TODO(), "tasks:queue").Result()
		reNum = tasknum
		idledevMap, _ := idledevices()
		DeviceNum := (*idledevMap)[p.Boardtype].Len()
		if DeviceNum > 0 {
			waitTime = reNum * 5
		} else {
			waitTime = reNum * 60
		}
	}
	if p.Queuetype == "groupqueue" {
		groupTasknum, _ := rdb.LLen(context.TODO(), "tasks:group:queue").Result()
		reNum = groupTasknum
		groupNum, _ := idlegroup(p.Boardtype)
		if groupNum > 0 {
			waitTime = reNum * 5
		} else {
			waitTime = reNum * 60
		}
	}
	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: waitTime})
}

func listuserimages(c *gin.Context) {
	userid := c.GetString("id")
	if len(userid) < 1 {
		err := errors.New("userid not exist")
		log.Error(err)
		c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
		return
	}

	// 命名空间
	namespace := tool.CreateMD5(userid)

	// 参数验证
	// p := request.UserImageName{}
	// if err := c.ShouldBindQuery(&p); err != nil {
	// 	err = fmt.Errorf("bing query parameter error {%v}", err)
	// 	log.Error(err)
	// 	c.SecureJSON(http.StatusOK, response.Response{Code: -1, Msg: err.Error()})
	// 	return
	// }
	usernamespace := fmt.Sprintf("%v", namespace)

	// 获取所有镜像列表
	registry, err := client.NewRegistry("http://"+dinfo.Registry.RegistryAddress, http.DefaultTransport)
	if err != nil {
		return
	}
	allImages := make([]string, 0)
	var last string
	for {
		images := make([]string, 1024)
		count, err := registry.Repositories(context.Background(), images, last)
		if err == io.EOF {
			allImages = append(allImages, images[:count]...)
			break
		} else if err != nil {
			return 
		}
		last = images[count-1]
		allImages = append(allImages, images...)
	}

	//获取用户镜像列表
	userImages := make([]string, 0)
	for i := 0; i < len(allImages); i++ {
		if usernamespace == strings.Split(allImages[i], "/")[0]{
			userImages = append(userImages, allImages[i])
		}
	}

	c.SecureJSON(http.StatusOK, response.Response{Code: 0, Msg: "success", Data: userImages})
	
}
