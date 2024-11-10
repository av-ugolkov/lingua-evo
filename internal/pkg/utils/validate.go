package utils

import (
	"regexp"
)

var (
	usernameRegex *regexp.Regexp
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
)

func init() {
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]{3,16}$`)
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9._\-]+\.[a-z]{2,4}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9.!@#$^&*(){}\[\]_%+\-]{8,20}$`)
}

func IsUsernameValid(u string) bool {
	return usernameRegex.MatchString(u)
}

func IsEmailValid(e string) bool {
	return emailRegex.MatchString(e)
}

func IsPasswordValid(p string) bool {
	return passwordRegex.MatchString(p)
}
