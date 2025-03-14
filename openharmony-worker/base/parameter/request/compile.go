package request

// CompileDownload 编译下载参数
type CompileDownload struct {
	CompileType string `form:"compiletype" binding:"required"`
	BoardType   string `form:"boardtype" binding:"required"`
	FileHash    string `form:"filehash" binding:"required"`
}

// CompileUpload 编译上传参数
type CompileUpload struct {
	CompileType string `form:"compiletype" json:"compiletype" binding:"required"`
	BoardType   string `form:"boardtype" json:"boardtype" binding:"required"`
	FileHash    string `form:"filehash" json:"filehash" binding:"required"`
}

// CompileSystemUpload 编译系统上传参数
type CompileSystemUpload struct {
	CompileType string `form:"compiletype" json:"compiletype" binding:"required"`
	BoardType   string `form:"boardtype" json:"boardtype" binding:"required"`
	FileHash    string `form:"filehash" json:"filehash" binding:"required"`
	Branch      string `form:"branch" json:"branch" binding:"required"`
}
