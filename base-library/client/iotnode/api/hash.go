package api

import (
	"crypto/sha256"
	"fmt"
)

// regenFileHash 重新生成文件哈希值
func (ad *Driver) regenFileHash(boardname string, oldfh string, groupid string, taskindex int) string {
	regenhash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%s:%v", boardname, oldfh, groupid, taskindex)))
	return fmt.Sprintf("%x", regenhash)
}
