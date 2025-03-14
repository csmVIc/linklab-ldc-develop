package api

// URLInfo url信息
type URLInfo struct {
	URL string `json:"url"`
}

// TokenInfo token信息
type TokenInfo struct {
	ChanSize int `json:"chansize"`
}

// AInfo 接口信息
type AInfo struct {
	FileDownload    URLInfo   `json:"filedownload"`
	PodYamlDownload URLInfo   `json:"podyamldownload"`
	TmpDir          string    `json:"tmpdir"`
	Token           TokenInfo `json:"token"`
}
