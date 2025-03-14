package table

// Tenant 用户表项
type Tenant struct {
	TenantID       int    `json:"tenantId" bson:"tenantId"`
	TenantName     string `json:"tenantName" bson:"tenantName"`
	IsSystemTenant bool   `json:"isSystemTenant" bson:"isSystemTenant"`
}

// TenantFilterByID 通过ID进行租户过滤
type TenantFilterByID struct {
	TenantID int `json:"tenantId" bson:"tenantId"`
}

// TenantFilterByFlag 通过权限进行租户
type TenantFilterByFlag struct {
	IsSystemTenant bool `json:"isSystemTenant" bson:"isSystemTenant"`
}
