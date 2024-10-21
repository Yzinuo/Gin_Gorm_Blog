package utils

import (
	"crypto/md5"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func BcrypyCheck(plian,hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(plian))
	return err == nil
}

func MD5(str string,b ...byte) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(b))
}