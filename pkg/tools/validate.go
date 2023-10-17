package tools

import (
	"regexp"
)

func IsEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func IsPasswordValid(p string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9.!@#$^&*(){}\[\]_%+\-]`)
	return emailRegex.MatchString(p)
}
