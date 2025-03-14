package fcache

import (
	"errors"
	"linklab/device-control-v2/base-library/auth"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	fcinfo *FCInfo
)

// RouterInit 设备操作路由初始化
func RouterInit(r *gin.Engine, fci *FCInfo) error {
	if fci == nil {
		err := errors.New("RouterInit di is nil")
		log.Error(err)
		return err
	}
	fcinfo = fci
	// 提供给用户的文件接口
	if err := userInit(r.Group("/api")); err != nil {
		log.Errorf("userInit error {%v}", err)
	}
	// 提供给客户端的文件接口
	if err := clientInit(r.Group("/api")); err != nil {
		log.Errorf("clientInit error {%v}", err)
	}
	return nil
}

func userInit(r *gin.RouterGroup) error {
	// 用户验证实例化
	handler, err := auth.New("users")
	if err != nil {
		log.Errorf("auth.New error {%v}", err)
		return err
	}
	// 烧写文件上传
	r.POST("/file", handler.AuthorizationRequired(), fileuploadhandler)
	// 运行日志下载
	r.GET("/devlog", handler.AuthorizationRequired(), devlogdownloadhandler)
	// PodYaml文件上传
	r.POST("/podyaml", handler.AuthorizationRequired(), podyamluploadhandler)
	// 镜像构建文件上传
	r.POST("/imagebuild", handler.AuthorizationRequired(), imagebuilduploadhandler)

	return nil
}

func clientInit(r *gin.RouterGroup) error {
	// 客户端验证实例化
	handler, err := auth.New("clients")
	if err != nil {
		log.Errorf("auth.New error {%v}", err)
		return err
	}
	// 烧写文件下载
	r.GET("/file", handler.AuthorizationRequired(), filedownloadhandler)
	// PodYaml文件下载
	r.GET("/podyaml", handler.AuthorizationRequired(), podyamldownloadhandler)
	// 镜像构建文件下载
	r.GET("/imagebuild", handler.AuthorizationRequired(), imagebuilddownloadhandler)
	return nil
}

func fileuploadhandler(c *gin.Context) {
	fileupload(c)
}

func filedownloadhandler(c *gin.Context) {
	filedownload(c)
}

func devlogdownloadhandler(c *gin.Context) {
	devlogdownload(c)
}

func podyamldownloadhandler(c *gin.Context) {
	podyamldownload(c)
}

func podyamluploadhandler(c *gin.Context) {
	podyamlupload(c)
}

func imagebuilddownloadhandler(c *gin.Context) {
	imagebuilddownload(c)
}

func imagebuilduploadhandler(c *gin.Context) {
	imagebuildupload(c)
}
