package ws

// WInfo websocket信息
type WInfo struct {
	TimeOut  int     `json:"timeout"`
	ChanSize int     `json:"chansize"`
	Msg      MsgInfo `json:"msg"`
}

// MsgInfo 消息队列的信息
type MsgInfo struct {
	Topic string `json:"topic"`
}
