package request

// RegisterParameter 注册参数
type RegisterParameter struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	TenantID int    `form:"tenantid" json:"tenantid"`
}
