package util

import "net/mail"

// IsValidEmail func to check if a string is a valid email address
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}
