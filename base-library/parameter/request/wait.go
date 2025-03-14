package request

// WaitTime 获取预计等待时长
type WaitTime struct {
	Queuetype string `form:"queuetype" json:"queuetype" binding:"required"`
	Boardtype string `form:"boardtype" json:"boardtype" binding:"required"`
}
