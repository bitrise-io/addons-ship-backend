package models

import "regexp"

// EmailVerifier ...
type EmailVerifier struct {
	Email string
}

// Verify ...
func (v EmailVerifier) Verify() bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !re.MatchString(v.Email) {
		return false
	}
	return true
}
