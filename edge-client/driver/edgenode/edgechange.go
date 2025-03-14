package edgenode

import (
	"fmt"
	"linklab/device-control-v2/base-library/parameter/msg"
	"linklab/device-control-v2/base-library/tool"

	log "github.com/sirupsen/logrus"
)

// GetEdgeNodesChange 获取设备列表变化
func (ed *Driver) GetEdgeNodesChange() (*msg.EdgeNodeStatusList, error) {

	// 加锁
	ed.edgenodes.edgeNodesLock.Lock()
	defer func() {
		ed.edgenodes.edgeNodesLock.Unlock()
	}()

	// 获取已有的设备列表
	tmpEdgeNodes := ed.getEdgeNodes()
	// 查询最新的设备列表
	nowEdgeNodes, err := ed.lsEdgeNodes()
	if err != nil {
		err = fmt.Errorf("ed.lsEdgeNodes error {%v}", err)
		log.Error(err)
		return nil, err
	}

	res := msg.EdgeNodeStatusList{
		HeartBeat: []msg.EdgeNodeStatus{},
		Delete:    []string{},
	}
	// 查询设备是否丢失
	for name := range *tmpEdgeNodes {
		if _, isOk := (*nowEdgeNodes)[name]; isOk == false {
			// 确认丢失
			log.Debugf("edgenode {%v} delete", name)
			ed.edgenodes.edgeNodesMap.Delete(name)
			res.Delete = append(res.Delete, name)
			continue
		}
	}
	// 查询设备是否新增
	for name := range *nowEdgeNodes {
		if _, isOk := (*tmpEdgeNodes)[name]; isOk == false {
			// 确认新增
			log.Debugf("edgenode {%v} add", name)
			// 标记边缘类型
			if err := ed.labelEdgeType(name); err != nil {
				err := fmt.Errorf("label edgenode error {%v}", err)
				log.Error(err)
				return nil, err
			}
			// 修改状态
			ed.edgenodes.edgeNodesMap.Store(name, (*nowEdgeNodes)[name])
			res.HeartBeat = append(res.HeartBeat, *(*nowEdgeNodes)[name])
			continue
		}

		// 查询设备状态是否发生变化
		if tool.CompareEdgeNodeStatus((*nowEdgeNodes)[name], (*tmpEdgeNodes)[name]) == false {
			// 确认变化
			value, isOk := ed.edgenodes.edgeNodesMap.Load(name)
			if isOk == false {
				err := fmt.Errorf("edgenode {%v} not in edgeNodesMap error", name)
				log.Error(err)
				return nil, err
			}
			edgenodestatus := value.(*msg.EdgeNodeStatus)
			edgenodestatus.Ready = (*(*nowEdgeNodes)[name]).Ready
			edgenodestatus.Architecture = (*(*nowEdgeNodes)[name]).Architecture
			edgenodestatus.OSImage = (*(*nowEdgeNodes)[name]).OSImage
			edgenodestatus.OS = (*(*nowEdgeNodes)[name]).OS
			edgenodestatus.Labels = (*(*nowEdgeNodes)[name]).Labels
			log.Debugf("edgenode {%v} status change to {%v}", name, *edgenodestatus)
			res.HeartBeat = append(res.HeartBeat, *(*nowEdgeNodes)[name])
		}
	}

	return &res, nil
}
