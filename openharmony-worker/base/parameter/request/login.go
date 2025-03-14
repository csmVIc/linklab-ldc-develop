package request

// LoginParameter 登录参数
type LoginParameter struct {
	ID       string `form:"id" json:"id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
