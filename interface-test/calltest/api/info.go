package api

// URLInfo url信息
type URLInfo struct {
	URL string `json:"url"`
}

// AInfo 接口信息
type AInfo struct {
	Login           URLInfo `json:"login"`
	FileUpload      URLInfo `json:"fileupload"`
	Burn            URLInfo `json:"burn"`
	WebSocket       URLInfo `json:"websocket"`
	CompileUpload   URLInfo `json:"compileupload"`
	CompileDownload URLInfo `json:"compiledownload"`
}
