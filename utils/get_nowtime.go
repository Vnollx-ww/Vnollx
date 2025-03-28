package utils

import (
	"log"
	"time"
)

func GetNowTime() (string, error) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println("获取时区失败:", err)
		return "", err
	}

	// 获取当前时间并转换为上海时区时间
	currentTime := time.Now().In(location)

	// 定义时间格式 (例如：yyyy-MM-dd HH:mm:ss)
	format := "2006-01-02 15:04:05"

	// 格式化上海时间
	formattedTime := currentTime.Format(format)
	return formattedTime, nil
}
