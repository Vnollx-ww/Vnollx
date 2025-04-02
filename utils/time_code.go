package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func toBase62(num int64) string {
	var result strings.Builder
	for num > 0 {
		result.WriteByte(base62[num%62])
		num /= 62
	}
	return result.String()
}
func GenerateUniqueString() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	base62Str := toBase62(timestamp)
	if len(base62Str) > 20 {
		base62Str = base62Str[:20]
	}
	rand.Seed(time.Now().UnixNano())
	randomSuffix := fmt.Sprintf("%02d", rand.Intn(100))
	result := base62Str + randomSuffix
	if len(result) > 20 {
		result = result[:20]
	}

	return result
}
