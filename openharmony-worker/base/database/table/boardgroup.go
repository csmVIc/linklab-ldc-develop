package table

// BoardGroup 开发板绑定组
type BoardGroup struct {
	Type   string   `json:"type" bson:"type"`
	Boards []string `json:"boards" bson:"boards"`
}

// BoardGroupFilter 开发板绑定组过滤器
type BoardGroupFilter struct {
	Type string `json:"type" bson:"type"`
}
