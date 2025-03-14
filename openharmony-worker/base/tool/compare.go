package tool

import "linklab/device-control-v2/base-library/parameter/msg"

func CompareEdgeNodeStatus(a *msg.EdgeNodeStatus, b *msg.EdgeNodeStatus) bool {

	if a == nil || b == nil {
		return false
	}

	if a.Ready != b.Ready ||
		a.Labels["linklab.edgetype"] != b.Labels["linklab.edgetype"] ||
		a.Labels["linklab.expand"] != b.Labels["linklab.expand"] {
		return false
	}

	return true
}
