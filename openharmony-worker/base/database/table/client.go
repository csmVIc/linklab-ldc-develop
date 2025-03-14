package table

// Client 客户端表项
type Client struct {
	UserName    string       `json:"username" bson:"username"`
	Password    string       `json:"password" bson:"password"`
	IsSuperUser bool         `json:"is_superuser" bson:"is_superuser"`
	Salt        string       `json:"salt" bson:"salt"`
	TenantID    map[int]bool `json:"tenantId" bson:"tenantId"`
}

// ClientFilter 客户端过滤
type ClientFilter struct {
	UserName string `json:"username" bson:"username"`
}
