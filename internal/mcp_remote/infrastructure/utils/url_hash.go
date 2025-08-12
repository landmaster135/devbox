package utils

import (
	"crypto/md5"
	"fmt"
)

// GenerateServerURLHash はサーバーURLからハッシュを生成する
func GenerateServerURLHash(serverURL string) string {
	hash := md5.Sum([]byte(serverURL))
	return fmt.Sprintf("%x", hash)
}
