package middlerware

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"golang.org/x/crypto/pbkdf2"
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
func HashPassword(password string, salt []byte) []byte {
	iterations := 100000 // 迭代次数
	keyLength := 64      // 生成密钥长度（字节）
	// 将字符串密码转换为字节
	passwordBytes := []byte(password)

	// 使用PBKDF2算法生成密钥
	return pbkdf2.Key(
		passwordBytes,
		salt,
		iterations,
		keyLength,
		sha512.New,
	)
}
func VerifyPassword(storedHash []byte, storedSalt []byte, inputPassword string) bool {
	inputHash := HashPassword(inputPassword, storedSalt)
	return subtle.ConstantTimeCompare(storedHash, inputHash) == 1
}
