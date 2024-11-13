package utils

import (
	"regexp"
)

var (
	nicknameRegex *regexp.Regexp
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
)

func init() {
	nicknameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,16}$`)
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9._\-]+\.[a-z]{2,4}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9.!@#$^&*(){}\[\]_%+\-]{8,20}$`)
}

func IsNicknameValid(u string) bool {
	return nicknameRegex.MatchString(u)
}

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}

func IsPasswordValid(p string) bool {
	return passwordRegex.MatchString(p)
}
