package table

import "time"

// DeviceLog 设备日志
type DeviceLog struct {
	UserID    string    `json:"userId" bson:"userId"`
	ClientID  string    `json:"clientId" bson:"clientId"`
	DevPort   string    `json:"devPort" bson:"devPort"`
	WaitingID string    `json:"waitingId" bson:"waitingId"`
	GroupID   string    `json:"groupId" bson:"groupId"`
	PID       string    `json:"pid" bson:"pid"`
	Logs      []string  `json:"logs" bson:"logs"`
	IsEnd     bool      `json:"isEnd" bson:"isEnd"`
	StartDate time.Time `json:"startDate" bson:"startDate"`
	EndDate   time.Time `json:"endDate" bson:"endDate"`
}

// DeviceLogFilter 设备日志过滤器
type DeviceLogFilter struct {
	WaitingID string `json:"waitingId" bson:"waitingId"`
}

// DeviceLogSetEnd 设备日志结束设置
type DeviceLogSetEnd struct {
	IsEnd   bool      `json:"isEnd" bson:"isEnd"`
	EndDate time.Time `json:"endDate" bson:"endDate"`
}
