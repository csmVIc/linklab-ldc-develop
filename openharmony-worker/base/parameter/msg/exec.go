package msg

// ExecErr 执行错误
type ExecErr struct {
	GroupID   string `form:"groupid" json:"groupid" binding:"required"`
	TaskIndex int    `form:"taskindex" json:"taskindex" binding:"required"`
	Msg       string `form:"msg" json:"msg" binding:"required"`
}
