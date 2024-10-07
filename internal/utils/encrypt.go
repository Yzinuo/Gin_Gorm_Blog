package utils

import "golang.org/x/crypto/bcrypt"

func BcrypyCheck(plian,hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(plian))
	return err == nil
}