package table

// User 用户表项
type User struct {
	UserID   string `json:"userId" bson:"userId"`
	Email    string `json:"email" bson:"email"`
	Hash     string `json:"hash" bson:"hash"`
	Salt     string `json:"salt" bson:"salt"`
	TenantID int    `json:"tenantId" bson:"tenantId"`
}

// UserFilter 用户过滤
type UserFilter struct {
	UserID string `json:"userId" bson:"userId"`
}
