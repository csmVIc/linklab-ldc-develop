package task

import (
	"linklab/device-control-v2/decision-maker/task/judge"
	"linklab/device-control-v2/decision-maker/task/pub"
	"linklab/device-control-v2/decision-maker/task/sub"
)

// TInfo 任务信息
type TInfo struct {
	Sub   sub.SInfo   `json:"sub"`
	Judge judge.JInfo `json:"judge"`
	Pub   pub.PInfo   `json:"pub"`
}
