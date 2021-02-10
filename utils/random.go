package utils

import (
	"math/rand"
	"time"
)

// 生成n位随机字符串
func RandCode(n int) string {
	code := make([]byte, n)
	rand.Seed(time.Now().Unix())
	dict := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	for i := 0; i < n; i++ {
		code[i] = dict[rand.Intn(len(dict)-1)]
	}
	return string(code)
}
