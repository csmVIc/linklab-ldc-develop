package response

// ClientStatus 客户端状态
type ClientStatus struct {
	ClientID string `json:"clientid"`
	TenantID []int  `json:"tenantid"`
}

// ClientList 客户端列表
type ClientList struct {
	Clients []ClientStatus `json:"clients"`
}
