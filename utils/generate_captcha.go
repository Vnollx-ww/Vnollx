package utils

import (
	"Vnollx/dal/redis"
	"context"
	"math/rand"
	"strconv"
	"time"
)

func GenerateCaptcha() string {
	rand.Seed(time.Now().UnixNano())
	captcha := rand.Intn(900000) + 100000
	return strconv.Itoa(captcha)
}
func GenerateUniqueCaptcha(ctx context.Context, phoneNumber string) (string, error) {
	var captcha string
	for {
		// 生成验证码
		captcha = GenerateCaptcha()
		// 检查验证码是否已经存在于 Redis
		exists, err := redis.IsCaptchaExists(ctx, captcha)
		if err != nil {
			return "", err
		}
		if !exists {
			break
		}
	}

	return captcha, nil
}
