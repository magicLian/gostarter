package util

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}
	suffixName := email[strings.Index(email, "@"):]
	if !strings.Contains(suffixName, ".") {
		return false
	}

	reg1, err1 := regexp.MatchString("^[a-z0-9_.-]+@[a-z0-9_-]+(.[a-z0-9_-]+)+", email)
	if err1 != nil || !reg1 {
		return false
	}

	if len(email) > 254 {
		return false
	}
	return true
}
