package problemsystem

import (
	"crypto/sha256"
	"fmt"
)

// ComputeWaitingID 创建等待ID
func ComputeWaitingID(groupid string, taskindex int) string {
	idseed := fmt.Sprintf("%s:%v", groupid, taskindex)
	idbinary := sha256.Sum256([]byte(idseed))
	return fmt.Sprintf("%x", idbinary)
}
