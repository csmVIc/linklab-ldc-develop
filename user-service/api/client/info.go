package client

// CInfo 客户端相关接口的配置信息
type CInfo struct {
	ClientTenant ClientTenantInfo `json:"clienttenant"`
}

// ClientTenantInfo 客户端和租户的配置信息
type ClientTenantInfo struct {
	ClientCacheTTL int `json:"clientcachettl"`
}
