package tool

import (
	"math/rand"
	"strings"
)

var randomchars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var randomstringlength = 16

// GenerateRandomString 生成随机字符串
func GenerateRandomString() string {
	var b strings.Builder
	for i := 0; i < randomstringlength; i++ {
		b.WriteRune(randomchars[rand.Intn(len(randomchars))])
	}
	return b.String()
}
