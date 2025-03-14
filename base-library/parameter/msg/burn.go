package msg

// BurnResult 烧写结果
type BurnResult struct {
	GroupID               string `form:"groupid" json:"groupid" binding:"required"`
	TaskIndex             int    `form:"taskindex" json:"taskindex" binding:"required"`
	Success               int    `form:"success" json:"success" binding:"required"`
	Msg                   string `form:"msg" json:"msg" binding:"required"`
	BeginBurnTime         int64  `form:"beginburntime" json:"beginburntime" binding:"required"`
	EndBurnTime           int64  `form:"endburntime" json:"endburntime" binding:"required"`
	BeginDownloadFileTime int64  `form:"begindownloadfiletime" json:"begindownloadfiletime" binding:"required"`
	EndDownloadFileTime   int64  `form:"enddownloadfiletime" json:"enddownloadfiletime" binding:"required"`
}
