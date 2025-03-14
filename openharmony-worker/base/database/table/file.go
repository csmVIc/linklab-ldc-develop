package table

import "go.mongodb.org/mongo-driver/bson/primitive"

// File 文件表项
type File struct {
	BoardName string           `json:"boardname" bson:"boardName"`
	FileHash  string           `json:"filehash" bson:"fileHash"`
	FileData  primitive.Binary `json:"filedata" bson:"fileData"`
}

// FileFilter 文件过滤
type FileFilter struct {
	BoardName string `json:"boardname" bson:"boardName"`
	FileHash  string `json:"filehash" bson:"fileHash"`
}
