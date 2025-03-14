package response

// BoardStatus 设备状态
type BoardStatus struct {
	BoardName string `json:"boardname"`
	BoardType string `json:"boardtype"`
}

// BoardList 设备列表
type BoardList struct {
	Boards []BoardStatus `json:"boards"`
}
