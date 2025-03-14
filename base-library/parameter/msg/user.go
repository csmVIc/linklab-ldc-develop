package msg

// UserMsgType 用户日志类型
type UserMsgType string

// TaskDataType 任务日志类型
type TaskDataType string

const (
	// UserControlMsg 用户控制报文
	UserControlMsg UserMsgType = "UserControlMsg"
	// TaskMsg 任务日志报文
	TaskMsg UserMsgType = "TaskMsg"
	// TaskAllocateMsg 任务分配结果
	TaskAllocateMsg TaskDataType = "TaskAllocateMsg"
	// TaskBurnMsg 任务烧写结果
	TaskBurnMsg TaskDataType = "TaskBurnMsg"
	// TaskExecErrMsg 任务执行错误结果
	TaskExecErrMsg TaskDataType = "TaskExecErrMsg"
	// TaskLogMsg 任务日志结果
	TaskLogMsg TaskDataType = "TaskLogMsg"
	// TaskEndRunMsg 任务设备结束运行
	TaskEndRunMsg TaskDataType = "TaskEndRunMsg"
)

// UserMsg 用户订阅的日志
type UserMsg struct {
	Code      int         `json:"code"`
	Type      UserMsgType `json:"type"`
	TimeStamp int64       `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// UserControlData 用户控制数据
// UserMsg.Type==UserControlMsg时，UserMsg.Data格式如下
type UserControlData struct {
	Msg string `json:"msg"`
}

// TaskData 任务日志数据
// UserMsg.Type==TaskMsg时，UserMsg.Data格式如下
type TaskData struct {
	GroupID   string       `json:"groupid"`
	TaskIndex int          `json:"taskindex"`
	Type      TaskDataType `json:"type"`
	Msg       string       `json:"msg"`
	Data      interface{}  `json:"data"`
}
