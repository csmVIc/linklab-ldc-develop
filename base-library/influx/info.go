package influx

// ClientInfo 客户端信息
type ClientInfo struct {
	URL                 string `json:"url"`
	UserName            string `json:"username"`
	PassWord            string `json:"password"`
	BatchSize           uint   `json:"batchsize"`
	FlushInterval       uint   `json:"flushinterval"`
	UseGZip             bool   `json:"usegzip"`
	DataBase            string `json:"database"`
	HealthCheckInterval int64  `json:"healthcheckinterval"`
}

// ChansInfo 等待队列信息
type ChansInfo struct {
	ThreadMultiple int   `json:"threadmultiple"`
	Size           int   `json:"size"`
	TimeOut        int64 `json:"timeout"`
}

// IInfo influxdb信息
type IInfo struct {
	Client ClientInfo `json:"client"`
	Chans  ChansInfo  `json:"chans"`
}

// Point Influx数据库的Point
type Point struct {
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
}
