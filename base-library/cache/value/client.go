package value

// ClientLoginStatus 客户端登录状态
type ClientLoginStatus struct {
	ClientID string       `json:"clientid"`
	TenantID map[int]bool `json:"tenantid"`
}
