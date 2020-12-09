package util

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// EncodeMd5 生成小写的md5值
func EncodeMd5(data string) string {
	src := md5.Sum([]byte(data))
	return hex.EncodeToString(src[:])
}

// EncodeMD5 生成大写的md5值
func EncodeMD5(data string) string {
	return strings.ToUpper(EncodeMd5(data))
}

// CreatePasswd 生成密码
func CreatePasswd(plainpwd, salt string) string {
	return EncodeMd5(plainpwd + salt)
}

// ValidatePasswd 验证密码正确性
func ValidatePasswd(plainpwd, salt, passwd string) bool {
	return EncodeMd5(plainpwd+salt) == passwd
}
