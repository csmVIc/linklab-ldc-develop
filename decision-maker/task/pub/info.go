package pub

// UserInfo 用户消息发布信息
type UserInfo struct {
	Topic string `json:"topic"`
}

// CmdBurnInfo 烧写命令信息
type CmdBurnInfo struct {
	Topic string `json:"topic"`
}

// ClientInfo 客户端消息发布信息
type ClientInfo struct {
	Burn CmdBurnInfo `json:"burn"`
}

// PInfo 消息发布信息
type PInfo struct {
	User   UserInfo   `json:"user"`
	Client ClientInfo `json:"client"`
}
