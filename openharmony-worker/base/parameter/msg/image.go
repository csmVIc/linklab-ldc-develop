package msg

// UserImageBuild 镜像构建
type UserImageBuild struct {
	UserID       string            `json:"userid"`
	FileHash     string            `json:"filehash"`
	ImageName    string            `json:"imagename"`
	NodeSelector map[string]string `json:"nodeselector"`
}
