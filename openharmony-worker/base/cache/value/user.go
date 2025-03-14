package value

// UserLoginStatus 用户登录状态
type UserLoginStatus struct {
	TrueCheck string `json:"truecheck"`
	Salt      string `json:"salt"`
	TenantID  int    `json:"tenantid"`
}
