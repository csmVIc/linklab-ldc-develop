package table

import "go.mongodb.org/mongo-driver/bson/primitive"

// PodYaml Pod配置文件
type PodYaml struct {
	FileHash string           `json:"filehash" bson:"fileHash"`
	FileData primitive.Binary `json:"filedata" bson:"fileData"`
}

// PodYamlFilter Pod配置文件过滤
type PodYamlFilter struct {
	FileHash string `json:"filehash" bson:"fileHash"`
}
