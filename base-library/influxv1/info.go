package influxv1

// ClientInfo 客户端信息
type ClientInfo struct {
	URL          string `json:"url"`
	UserName     string `json:"username"`
	PassWord     string `json:"password"`
	DataBase     string `json:"database"`
	WriteTimeOut int    `json:"writetimeout"`
	PingTimeOut  int    `json:"pingtimeout"`
}

// IInfo 初始化信息
type IInfo struct {
	Client ClientInfo `json:"client"`
}
