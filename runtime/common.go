package runtime

import (
	"math/rand"
	"strings"
	"time"
)

const (
	EmptyString = ""
	EmptyJson   = "{}"
)

func GenerateNickname() string {
	return generateRandomString(7)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var sb strings.Builder

	for i := 0; i < length; i++ {
		randomIndex := seededRand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
