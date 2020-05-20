package salt

import (
	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd string) (res string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)

	return string(hash), err
}

func ComparePasswords(hashedPwd, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePlain := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)
	return err == nil
}
