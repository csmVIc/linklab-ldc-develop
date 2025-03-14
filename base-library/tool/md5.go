package tool

import (
	"crypto/md5"
	"fmt"
)

func CreateMD5(value string) string {
	hash := md5.Sum([]byte(value))
	return fmt.Sprintf("%x", hash)
}
