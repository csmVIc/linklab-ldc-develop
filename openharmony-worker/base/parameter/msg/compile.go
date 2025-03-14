package msg

// CompileTask 编译任务
type CompileTask struct {
	CompileType string `json:"compiletype"`
	BoardType   string `json:"boardtype"`
	FileHash    string `json:"filehash"`
}
