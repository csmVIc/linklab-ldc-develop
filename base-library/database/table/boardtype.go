package table

// BoardType 系统支持开发板类型
type BoardType struct {
	BoardType string `json:"boardtype" bson:"boardType"`
	AllowCmd  bool   `json:"allowcmd" bson:"allowCmd"`
}

// BoardTypeFilter 开发板类型过滤
type BoardTypeFilter struct {
	BoardType string `json:"boardtype" bson:"boardType"`
}
