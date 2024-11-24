package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func EncryptAES(plainText, key string) (string, error) {
	keyBytes := []byte(key)
	plainTextBytes := []byte(plainText)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("pkg.aes.EncryptAES: %w", err)
	}

	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("pkg.aes.EncryptAES: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	return hex.EncodeToString(cipherText), nil
}

func DecryptAES(cipherHex, key string) (string, error) {
	keyBytes := []byte(key)
	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", fmt.Errorf("pkg.aes.DecryptAES: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("pkg.aes.DecryptAES: %w", err)
	}

	if len(cipherText) < aes.BlockSize {
		return "", fmt.Errorf("pkg.aes.DecryptAES: ciphertext too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
