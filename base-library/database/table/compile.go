package table

import "go.mongodb.org/mongo-driver/bson/primitive"

// CompileTable 编译数据表
type CompileTable struct {
	Type        string           `json:"type" bson:"type"`
	CompileType string           `json:"compileType" bson:"compileType"`
	BoardType   string           `json:"boardType" bson:"boardType"`
	FileHash    string           `json:"fileHash" bson:"fileHash"`
	Branch      string           `json:"branch" bson:"branch"`
	FileData    primitive.Binary `json:"fileData" bson:"fileData"`
	Output      primitive.Binary `json:"output" bson:"output"`
	Message     string           `json:"message" bson:"message"`
	Status      string           `json:"status" bson:"status"`
}

// CompileTableFilter 编译数据表过滤器
type CompileTableFilter struct {
	CompileType string `json:"compileType" bson:"compileType"`
	BoardType   string `json:"boardType" bson:"boardType"`
	FileHash    string `json:"fileHash" bson:"fileHash"`
}
