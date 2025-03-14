package table

// Board 系统支持开发板的表项
type Board struct {
	BoardName string `json:"boardname" bson:"boardName"`
	BoardType string `json:"boardtype" bson:"boardType"`
}

// BoardFilter 开发板过滤
type BoardFilter struct {
	BoardName string `json:"boardname" bson:"boardName"`
}
