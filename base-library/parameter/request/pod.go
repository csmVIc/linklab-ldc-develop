package request

import (
	"fmt"
	"strings"
)

// PodYamlDownload PodYaml文件下载
type PodYamlDownload struct {
	FileHash string `form:"filehash" json:"filehash" binding:"required"`
}

// EdgeClientPodApply Pod部署参数
type EdgeClientPodApply struct {
	Namespace       string `form:"namespace" json:"namespace" binding:"required"`
	YamlHash        string `form:"yamlhash" json:"yamlhash" binding:"required"`
	UseEdgeRegistry bool   `form:"useedgeregistry" json:"useedgeregistry"`
	CreateIngress   bool   `form:"createingress" json:"createingress"`
	NodeAddSelector string `form:"nodeaddselector" json:"nodeaddselector"`
}

// EdgeClientPodDelete Pod删除参数
type EdgeClientPodDelete struct {
	Namespace string `form:"namespace" json:"namespace" binding:"required"`
	Pod       string `form:"pod" json:"pod" binding:"required"`
}

// EdgeClientPodLog Pod日志参数
type EdgeClientPodLog struct {
	Namespace string `form:"namespace" json:"namespace" binding:"required"`
	Pod       string `form:"pod" json:"pod" binding:"required"`
	Container string `form:"container" json:"container"`
}

// QueryRaw 序列化为字符串
func (pl *EdgeClientPodLog) QueryRaw() string {
	queryraw := fmt.Sprintf("namespace=%v&pod=%v", pl.Namespace, pl.Pod)
	if len(pl.Container) > 0 {
		queryraw = fmt.Sprintf("%v&container=%v", queryraw, pl.Container)
	}
	return queryraw
}

// EdgeClientPodExec Pod执行参数
type EdgeClientPodExec struct {
	Namespace string   `form:"namespace" json:"namespace" binding:"required"`
	Pod       string   `form:"pod" json:"pod" binding:"required"`
	Container string   `form:"container" json:"container"`
	Commands  []string `form:"commands[]" json:"commands" binding:"required"`
}

// QueryRaw 序列化为字符串
func (pe *EdgeClientPodExec) QueryRaw() string {
	queryraw := fmt.Sprintf("namespace=%v&pod=%v", pe.Namespace, pe.Pod)
	if len(pe.Container) > 0 {
		queryraw = fmt.Sprintf("%v&container=%v", queryraw, pe.Container)
	}
	sb := &strings.Builder{}
	for i, cmd := range pe.Commands {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString("commands[]=")
		sb.WriteString(cmd)
	}
	queryraw = fmt.Sprintf("%v&%v", queryraw, sb.String())
	return queryraw
}

// UserPodApply Pod部署参数
type UserPodApply struct {
	YamlHash        string `form:"yamlhash" json:"yamlhash" binding:"required"`
	UseEdgeRegistry bool   `form:"useedgeregistry" json:"useedgeregistry"`
	CreateIngress   bool   `form:"createingress" json:"createingress"`
}

// UserPodLog Pod日志参数
type UserPodLog struct {
	Pod       string `form:"pod" json:"pod" binding:"required"`
	Container string `form:"container" json:"container"`
	ClientID  string `form:"clientid" json:"clientid" binding:"required"`
}

// QueryRaw 序列化为字符串
func (pl *UserPodLog) QueryRaw() string {
	queryraw := fmt.Sprintf("clientid=%v&pod=%v", pl.ClientID, pl.Pod)
	if len(pl.Container) > 0 {
		queryraw = fmt.Sprintf("%v&container=%v", queryraw, pl.Container)
	}
	return queryraw
}

// UserPodQuery Pod列表
type UserPodQuery struct {
	AllPods bool `form:"allpods" json:"allpods"`
}

// UserPodDelete Pod删除参数
type UserPodDelete struct {
	ClientID string `form:"clientid" json:"clientid" binding:"required"`
	Pod      string `form:"pod" json:"pod" binding:"required"`
}

// UserPodExec Pod执行参数
type UserPodExec struct {
	ClientID  string   `form:"clientid" json:"clientid" binding:"required"`
	Pod       string   `form:"pod" json:"pod" binding:"required"`
	Container string   `form:"container" json:"container"`
	Commands  []string `form:"commands[]" json:"commands" binding:"required"`
}

// QueryRaw 序列化为字符串
func (pe *UserPodExec) QueryRaw() string {
	queryraw := fmt.Sprintf("clientid=%v&pod=%v", pe.ClientID, pe.Pod)
	if len(pe.Container) > 0 {
		queryraw = fmt.Sprintf("%v&container=%v", queryraw, pe.Container)
	}
	sb := &strings.Builder{}
	for i, cmd := range pe.Commands {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString("commands[]=")
		sb.WriteString(cmd)
	}
	queryraw = fmt.Sprintf("%v&%v", queryraw, sb.String())
	return queryraw
}
