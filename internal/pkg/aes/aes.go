package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
)

type Cipher struct {
	key   string
	block cipher.Block
}

var refreshToken *Cipher
var mx sync.Mutex

func EncryptAES(plainText, key string) (string, error) {
	plainTextBytes := []byte(plainText)

	block, err := cipherBlock(key)
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
	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", fmt.Errorf("pkg.aes.DecryptAES: %w", err)
	}

	block, err := cipherBlock(key)
	if err != nil {
		return "", fmt.Errorf("pkg.aes.EncryptAES: %w", err)
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

func cipherBlock(key string) (cipher.Block, error) {
	if refreshToken == nil {
		mx.Lock()
		defer mx.Unlock()
	}
	if refreshToken != nil && refreshToken.key == key {
		return refreshToken.block, nil
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	refreshToken = &Cipher{
		key:   key,
		block: block,
	}

	return refreshToken.block, nil
}
