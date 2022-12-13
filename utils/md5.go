package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

//md5注册加密转小写

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	//转换md5
	temp := h.Sum(nil)
	return hex.EncodeToString(temp)
}

//大写

func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

//加密

func MakePasssword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 解密

func ValidPasssword(plainpwd, salt string, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
