package logger

// LogInfo 单个日志记录器信息
type LogInfo struct {
	Pointname string `json:"pointname"`
	Directory string `json:"directory"`
}

// LInfo 日志记录信息
type LInfo struct {
	Logs map[string]LogInfo `json:"logs"`
}
