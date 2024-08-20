package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckHashPassword(hash, password string)  error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}


