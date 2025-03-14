package device

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	dinfo *DInfo
)

// RouterInit 设备操作路由初始化
func RouterInit(r *gin.RouterGroup, di *DInfo) error {
	if di == nil {
		err := errors.New("RouterInit di is nil")
		log.Error(err)
		return err
	}
	dinfo = di
	// 设备程序烧写，设备状态列表，设备支持列表
	r.POST("/device/burn", burnhandler)
	r.GET("/device/list", listhandler)
	r.GET("/board/list", boardhandler)

	r.GET("/device/listuserdevice", listuserdevicehandler)
	r.POST("/device/cmd", cmdhandler)

	r.POST("/device/creategroup", creategrouphandler)
	r.POST("/device/allocategroup", allocategrouphandler)
	r.POST("/device/linkgroup", linkgrouphandler)
	r.POST("/device/unlinkgroup", unlinkgrouphandler)
	r.GET("/device/listlinkgroup", listlinkgrouphandler)
	r.GET("/device/listdefinegroup", listdefinegrouphandler)

	// 获取预计等待时间
	r.GET("/device/getcurrtasknum", tasknumhandler)

	// 获取用户镜像列表
	r.GET("/device/listuserimages", listuserimageshandler)
	return nil
}

func burnhandler(c *gin.Context) {
	deviceburn(c, dinfo.TaskNumLimit.MinTaskNum, dinfo.TaskNumLimit.MaxTaskNum, dinfo.TaskRuntimeLimit.MinRuntime, dinfo.TaskRuntimeLimit.MaxRuntime, dinfo.Msg.Topic, dinfo.Msg.ReplyTimeOut)
}

func listhandler(c *gin.Context) {
	listalldevice(c)
}

func boardhandler(c *gin.Context) {
	boardlist(c)
}

func listuserdevicehandler(c *gin.Context) {
	listuserdevice(c)
}

func cmdhandler(c *gin.Context) {
	sendcmd(c, dinfo.Cmd.Topic, dinfo.Cmd.ReplyTimeOut)
}

func creategrouphandler(c *gin.Context) {
	creategroup(c)
}

func allocategrouphandler(c *gin.Context) {
	allocategroup(c)
}

func linkgrouphandler(c *gin.Context) {
	linkgroup(c)
}

func listlinkgrouphandler(c *gin.Context) {
	listlinkgroup(c)
}

func listdefinegrouphandler(c *gin.Context) {
	listdefinegroup(c)
}

func unlinkgrouphandler(c *gin.Context) {
	unlinkgroup(c)
}

func tasknumhandler(c *gin.Context){
	getcurrtasknum(c)
}

func listuserimageshandler(c *gin.Context){
	listuserimages(c)
}
