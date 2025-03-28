package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateFileHash(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 使用 SHA-256 算法计算哈希
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	// 返回文件的哈希值
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
