package utils

import (
	"time"
)

func GetNowTime() (string, error) {
	// 手动设置时区为上海时区（UTC+8）
	shanghai := time.FixedZone("Asia/Shanghai", 8*60*60)

	// 获取当前时间并转换为上海时区时间
	currentTime := time.Now().In(shanghai)

	// 定义时间格式 (例如：yyyy-MM-dd HH:mm:ss)
	format := "2006-01-02 15:04:05"

	// 格式化上海时间
	formattedTime := currentTime.Format(format)
	return formattedTime, nil
}
