package user

// IoTNodeInfo 开发板节点
type IoTNodeInfo struct {
	ListURL string `json:"listurl"`
	CmdURL  string `json:"cmdurl"`
}

// EdgeNodeInfo 边缘节点
type EdgeNodeInfo struct {
	NodeListURL          string `json:"nodelisturl"`
	NodeResourceURL      string `json:"noderesourceurl"`
	PodListURL           string `json:"podlisturl"`
	PodResourceURL       string `json:"podresourceurl"`
	PodYamlUploadURL     string `json:"podyamluploadurl"`
	PodApplyURL          string `json:"podapplyurl"`
	PodDeleteURL         string `json:"poddeleteurl"`
	PodLogURL            string `json:"podlogurl"`
	PodExecURL           string `json:"podexecurl"`
	ImageSourceUploadURL string `json:"imagesourceuploadurl"`
	ImageBuildURL        string `json:"imagebuildurl"`
}

// UserInfo 用户信息
type UserInfo struct {
	LoginURL string `json:"loginurl"`
}

// DeviceBindGroupInfo 设备绑定组信息
type DeviceBindGroupInfo struct {
	CreateGroupURL     string `json:"creategroupurl"`
	LinkGroupURL       string `json:"linkgroupurl"`
	UnLinkGroupURL     string `json:"unlinkgroupurl"`
	ListLinkGroupURL   string `json:"listlinkgroupurl"`
	ListDefineGroupURL string `json:"listdefinegroupurl"`
}

// UInfo 用户接口信息
type UInfo struct {
	IoTNode         IoTNodeInfo         `json:"iotnode"`
	EdgeNode        EdgeNodeInfo        `json:"edgenode"`
	User            UserInfo            `json:"user"`
	DeviceBindGroup DeviceBindGroupInfo `json:"devicebindgroup"`
}
