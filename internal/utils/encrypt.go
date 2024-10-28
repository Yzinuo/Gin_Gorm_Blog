package utils

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// 密码存储再数据库时，肯定不能明文存储，需要加密
func BcryptHash (str string) (string, error){
	bytes,err := bcrypt.GenerateFromPassword([]byte(str),bcrypt.DefaultCost)
	return string(bytes),err
}

// 使用bcrypt 对比明文字符串和加密后的哈希字符串
func BcrypyCheck(plian,hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(plian))
	return err == nil
}

func MD5(str string,b ...byte) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(b))
}