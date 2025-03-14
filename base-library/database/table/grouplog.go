package table

// GroupLog 组日志
type GroupLog struct {
	GroupID    string   `json:"groupId" bson:"groupId"`
	UserID     string   `json:"userId" bson:"userId"`
	WaitingIDs []string `json:"waitingIds" bson:"waitingIds"`
	PID        string   `json:"pid" bson:"pid"`
	Logs       []string `json:"logs" bson:"logs"`
}

// GroupLogFilter 组日志过滤器
type GroupLogFilter struct {
	GroupID string `json:"groupId" bson:"groupId"`
}
