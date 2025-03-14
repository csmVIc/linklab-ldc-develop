package edgenode

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	einfo *EInfo
)

// RouterInit 路由初始化
func RouterInit(r *gin.RouterGroup, ei *EInfo) error {
	if ei == nil {
		err := errors.New("RouterInit ei is nil")
		log.Error(err)
		return err
	}
	einfo = ei

	// 边缘节点查询
	r.GET("/edgenode/list", listnodehandler)
	r.GET("/pod/list", listpodhandler)
	r.GET("/edgenode/resource/list", listnoderesourcehandler)
	r.GET("/pod/resource/list", listpodresourcehandler)
	r.POST("/pod/apply", podapplyhandler)
	r.GET("/pod/log", podloghandler)
	r.POST("/pod/delete", poddeletehandler)
	r.GET("/pod/exec", podexechandler)
	r.POST("/image/build", imagebuildhandler)
	return nil
}

func listnodehandler(c *gin.Context) {
	listnode(c)
}

func listpodhandler(c *gin.Context) {
	listpod(c)
}

func listnoderesourcehandler(c *gin.Context) {
	listnoderesource(c)
}

func listpodresourcehandler(c *gin.Context) {
	listpodresource(c)
}

func podapplyhandler(c *gin.Context) {
	podapply(c)
}

func podloghandler(c *gin.Context) {
	podlog(c)
}

func poddeletehandler(c *gin.Context) {
	poddelete(c)
}

func podexechandler(c *gin.Context) {
	podexec(c)
}

func imagebuildhandler(c *gin.Context) {
	imagebuild(c)
}
