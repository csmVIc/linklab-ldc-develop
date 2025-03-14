package request

// FileUpload 烧写文件上传
type FileUpload struct {
	BoardName string `form:"boardname" json:"boardname" binding:"required"`
}

// FileDownload 烧写文件下载
type FileDownload struct {
	BoardName string `form:"boardname" json:"boardname" binding:"required"`
	FileHash  string `form:"filehash" json:"filehash" binding:"required"`
}
