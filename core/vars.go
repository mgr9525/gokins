package core

import (
	"math/rand"
	"time"
)

var (
	JwtKey        = ""
	JwtCookiePath = "/"
)

// RandomString 随机生成字符串
func RandomString(l int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	bts := []byte(str)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bts[r.Intn(len(bts))])
	}
	return string(result)
}
