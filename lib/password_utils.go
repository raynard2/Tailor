package lib

import (
	"golang.org/x/crypto/bcrypt"
)

func CreateHashFromPassword(password string) (string,error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		panic("error generating hash from password")
		}

	return string(hash), err
}


func CompareHashAndPassword (hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password))
	return err==nil
}
