package auth

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

// CreateToken 创建验证秘钥
func CreateToken(id string, hash string) string {
	ntime := time.Now().Unix()
	rint := rand.Int()
	tokenseed := fmt.Sprintf("%s:%s:%v:%v", id, hash, ntime, rint)
	tokenbinary := sha256.Sum256([]byte(tokenseed))
	return fmt.Sprintf("%x", tokenbinary)
}
