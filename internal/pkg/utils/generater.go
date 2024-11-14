package utils

import (
	"math/rand/v2"

	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordSolt = 11
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordSolt)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GenerateCode() int {
	return rand.IntN(999999-100000) + 100000
}
