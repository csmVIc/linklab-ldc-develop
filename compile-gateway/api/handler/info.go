package handler

// HInfo 接口的配置参数
type HInfo struct {
	Download      DownloadInfo      `json:"download"`
	Upload        UploadInfo        `json:"upload"`
	CompileSystem CompileSystemInfo `json:"compilesystem"`
}

// DownloadInfo 下载接口的参数
type DownloadInfo struct {
	Timeout int64 `json:"timeout"`
}

// UploadInfo 上传接口的参数
type UploadInfo struct {
	ReplyTimeout int64  `json:"replytimeout"`
	Topic        string `json:"topic"`
}

// CompileSystemInfo 上传编译系统的参数
type CompileSystemInfo struct {
	ReplyTimeout int64  `json:"replytimeout"`
	Topic        string `json:"topic"`
}
