package utils

/*
* 密码的加密和验证
 */
import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword 根据普通密码生成加密后的密码
func GeneratePassword(Password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword 比较加密后的密码和加密前的密码，进行验证
func VerifyPassword(Password string, Hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(Hash), []byte(Password)) //验证（对比）
	return err
}
