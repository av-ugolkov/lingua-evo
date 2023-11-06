package tools

import (
	"crypto/sha256"
	"fmt"

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

func HashValue(value string) string {
	bytes := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", bytes)
}
