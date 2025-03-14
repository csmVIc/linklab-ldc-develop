package user

// UInfo 用户信息
type UInfo struct {
	Login    LoginInfo    `json:"login"`
	Register RegisterInfo `json:"register"`
}

// LoginInfo 登录参数
type LoginInfo struct {
	TimeOut int `json:"timeout"`
}

// RegisterInfo 注册参数
type RegisterInfo struct {
	Email string `json:"email"`
}
