package problemsystem

// PInfo 判题接口
type PInfo struct {
	URL string `json:"url"`
}

// PResult 判题接口结果
type PResult struct {
	Result string `json:"result"`
	Msg    string `json:"msg"`
}

// PrepareRequest 判题准备接口
type PrepareRequest struct {
	WaitingID string `json:"waitingId"`
	PID       string `json:"pid"`
}

// StartRequest 判题开始接口
type StartRequest struct {
	WaitingID string `json:"waitingId"`
	PID       string `json:"pid"`
	Email     bool   `json:"email"`
}
