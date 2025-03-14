package compile

// Command 编译命令执行相关信息
type Command struct {
	Type    string `json:"type"`
	SupSys  bool   `json:"supsys"`
	Branch  string `json:"branch"`
	Rootdir string `json:"rootdir"`
	Indir   string `json:"indir"`
	Outdir  string `json:"outdir"`
	Rregex  string `json:"rregex"`
	Cmd     string `json:"cmd"`
	ErrFlag string `json:"errflag"`
}

// DirectoryInfo 目录信息
type DirectoryInfo struct {
	Tmp       string `json:"tmp"`
	Workspace string `json:"workspace"`
	InitZip   string `json:"initzip"`
}

// ChannelInfo 队列信息
type ChannelInfo struct {
	Size    int `json:"size"`
	TimeOut int `json:"timeout"`
}

// CInfo 编译信息
type CInfo struct {
	Directory DirectoryInfo      `json:"directory"`
	Commands  map[string]Command `json:"commands"`
	Channel   ChannelInfo        `json:"channel"`
}
