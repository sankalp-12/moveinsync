package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(passwordDB string, passwordRequest string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordDB), []byte(passwordRequest))
	return err
}
