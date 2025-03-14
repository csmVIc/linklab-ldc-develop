package fcache

// SizeRangeInfo 大小范围
type SizeRangeInfo struct {
	MaxBytes int `json:"maxbytes"`
	MinBytes int `json:"minbytes"`
}

// FCInfo 文件缓存信息
type FCInfo struct {
	SizeRange SizeRangeInfo `json:"sizerange"`
}
