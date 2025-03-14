package monitor

// MetricsInfo 监控信息
type MetricsInfo struct {
	Interval  int64  `json:"interval"`
	Namespace string `json:"namespace"`
}

// MInfo 监控器信息
type MInfo struct {
	MetricsMap map[string]MetricsInfo `json:"metricsmap"`
}
