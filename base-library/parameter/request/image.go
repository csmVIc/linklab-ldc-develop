package request

// EdgeClientImageBuild 镜像打包参数
type EdgeClientImageBuild struct {
	Namespace    string            `form:"namespace" json:"namespace" binding:"required"`
	FileHash     string            `form:"filehash" json:"filehash" binding:"required"`
	ImageName    string            `form:"imagename" json:"imagename" binding:"required"`
	NodeSelector map[string]string `form:"nodeselector" json:"nodeselector"`
}

// ImageBuildDownload 镜像构建文件下载
type ImageBuildDownload struct {
	FileHash string `form:"filehash" json:"filehash" binding:"required"`
}

// UserImageBuild 镜像构建参数
type UserImageBuild struct {
	FileHash     string            `form:"filehash" json:"filehash" binding:"required"`
	ImageName    string            `form:"imagename" json:"imagename" binding:"required"`
	NodeSelector map[string]string `form:"nodeselector" json:"nodeselector"`
}

type UserImageName struct{
	ImageName    string            `form:"imagename" json:"imagename" binding:"required"`
}