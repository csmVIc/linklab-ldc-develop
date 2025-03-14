package msg

// ReplyMsg 回复报文
type ReplyMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
