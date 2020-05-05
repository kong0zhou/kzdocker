package utils

import (
	"fmt"
	"iron/log"

	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword encrypt the password,it is irreversible
func EncryptPassword(password string) (ciphertext string, err error) {
	if password == `` {
		err = fmt.Errorf(`密码为空，但这显然不可能，你前面的代码肯定有问题`)
		log.Error(err.Error())
		return ``, err
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Error(err.Error())
		return ``, err
	}
	return string(bytes), nil
}

// CheckPassword check password
func CheckPassword(ciphertext string, password string) bool {
	if ciphertext == "" || password == "" {
		log.Error(`ciphertext or password is empty`)
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(ciphertext), []byte(password))
	return err == nil
}
