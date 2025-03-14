package request

// type ClientListQuery struct {
// 	ClientID string `form:"clientid" json:"clientid" binding:"required"`
// }

// ClientTenantSet 配置客户端所属租户
type ClientTenantSet struct {
	ClientID string `form:"clientid" json:"clientid" binding:"required"`
	TenantID int    `form:"tenantid" json:"tenantid" binding:"required"`
}

// ClientTenantChange 修改客户端所属租户
type ClientTenantChange struct {
	ClientID       string `form:"clientid" json:"clientid" binding:"required"`
	SourceTenantID int    `form:"sourcetenantid" json:"sourcetenantid" binding:"required"`
	DestTenantID   int    `form:"desttenantid" json:"desttenantid" binding:"required"`
}
