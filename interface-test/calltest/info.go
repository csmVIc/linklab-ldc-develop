package calltest

import (
	"linklab/device-control-v2/interface-test/calltest/api"
)

// LoginInfo 登录信息
type LoginInfo struct {
	UserName string `json:"username"`
	PassWord string `json:"password"`
}

// CompileInfo 编译信息
type CompileInfo struct {
	NeedCompile    bool   `json:"needcompile"`
	SourcePath     string `json:"sourcepath"`
	RandomFileName string `json:"randomfilename"`
	CompileType    string `json:"compiletype"`
	BoardType      string `json:"boardtype"`
}

// EndChanInfo 结束通道信息
type EndChanInfo struct {
	Size int `json:"size"`
}

// GroupInfo 测试组信息
type GroupInfo struct {
	BeginID    int `json:"beginid"`
	TotalTimes int `json:"totaltimes"`
	Step       int `json:"step"`
	EndTimes   int `json:"endtimes"`
}

// TestInfo 测试信息
type TestInfo struct {
	Groups  []GroupInfo `json:"groups"`
	DataDir string      `json:"datadir"`
}

// BurnInfo 烧写信息
type BurnInfo struct {
	BoardName  string `json:"boardname"`
	RunTime    int    `json:"runtime"`
	FileRandom bool   `json:"filerandom"`
	FilePath   string `json:"filepath"`
	FileSize   int    `json:"filesize"`
	RetryTimes int    `json:"retrytimes"`
}

// OSSQueryInfo 数据库查询信息
type OSSQueryInfo struct {
	PodMetrics  string `json:"podmetrics"`
	NodeMetrics string `json:"nodemetrics"`
}

// CInfo 接口测试
type CInfo struct {
	Login    LoginInfo    `json:"login"`
	Test     TestInfo     `json:"test"`
	Burn     BurnInfo     `json:"burn"`
	EndChan  EndChanInfo  `json:"endchan"`
	API      api.AInfo    `json:"api"`
	OSSQuery OSSQueryInfo `json:"ossquery"`
	Compile  CompileInfo  `json:"compile"`
}

// TestThreadEnd 测试单例结束
type TestThreadEnd struct {
	UserName string
	Err      error
}

// TaskTestResult 任务测试结果
type TaskTestResult struct {
	AvgLogRecvResp       int    `json:"avglogrecvresp"`
	AvgDeviceLogRecvResp int    `json:"avgdevicelogrecvresp"`
	AvgNatsLogRecvResp   int    `json:"avgnatslogrecvresp"`
	AvgUserLogRecvResp   int    `json:"avguserlogrecvresp"`
	BeginDeviceBurn      int64  `json:"begindeviceburn"`
	BeginFileUpload      int64  `json:"beginfileupload"`
	BeginUserLogin       int64  `json:"beginuserlogin"`
	DeviceAllocate       int64  `json:"deviceallocate"`
	EndDeviceBurn        int64  `json:"enddeviceburn"`
	EndDeviceRun         int64  `json:"enddevicerun"`
	EndFileUpload        int64  `json:"endfileupload"`
	EndUserLogin         int64  `json:"enduserlogin"`
	EnterTasksWait       int64  `json:"entertaskswait"`
	GroupID              string `json:"groupid"`
	LogBytesCount        int    `json:"logbytescount"`
	LogRecvCount         int    `json:"logrecvcount"`
	UserName             string `json:"username"`
	BeginCompile         int64  `json:"begincompile"`
	EndCompile           int64  `json:"endcompile"`
	WebSocketBroke       bool   `json:"websocketbroke"`
}
