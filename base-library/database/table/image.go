package table

import "go.mongodb.org/mongo-driver/bson/primitive"

// ImageBuild 镜像构建文件
type ImageBuild struct {
	FileHash string           `json:"filehash" bson:"fileHash"`
	FileData primitive.Binary `json:"filedata" bson:"fileData"`
}

// ImageBuildFilter 镜像构建文件过滤
type ImageBuildFilter struct {
	FileHash string `json:"filehash" bson:"fileHash"`
}
