package cmd

// UserLoginInfo 用户登录信息
type UserLoginInfo struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

var userlogininfo UserLoginInfo
