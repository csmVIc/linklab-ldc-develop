package auth

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

// CreateSalt 创建用户盐值
func CreateSalt(id string) string {
	ntime := time.Now().Unix()
	rint := rand.Int()
	saltseed := fmt.Sprintf("%s:%v:%v", id, ntime, rint)
	saltbinary := sha256.Sum256([]byte(saltseed))
	return fmt.Sprintf("%x", saltbinary)
}
