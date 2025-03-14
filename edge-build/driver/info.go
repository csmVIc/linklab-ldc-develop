package driver

// APIInfo 接口信息
type APIInfo struct {
	BuildDownloadURL string `json:"builddownloadurl"`
	FileHash         string `json:"filehash"`
	Token            string `json:"token"`
}

// BuildInfo 构建信息
type BuildInfo struct {
	ImageName string `json:"imagename"`
}

// DInfo 配置信息
type DInfo struct {
	API   APIInfo   `json:"api"`
	Build BuildInfo `json:"build"`
}
