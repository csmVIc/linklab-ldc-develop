package api

// TokenUpdate 秘钥更新
func (ad *Driver) TokenUpdate() {
	for newtoken := range ad.tokenchan {
		ad.token = newtoken
	}
}

// GetTokenChan 获取token通道
func (ad *Driver) GetTokenChan() *chan string {
	return &ad.tokenchan
}

// GetToken 获取token
func (ad *Driver) GetToken() string {
	return ad.token
}
