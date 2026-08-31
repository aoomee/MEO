package auth

import (
	"errors"
	"unicode/utf8"
)

// ValidateNewPassword applies one policy at every password creation/reset
// boundary. The 72-byte ceiling is imposed by bcrypt itself.
func ValidateNewPassword(password string) error {
	if utf8.RuneCountInString(password) < 12 {
		return errors.New("password must contain at least 12 characters")
	}
	if len([]byte(password)) > 72 {
		return errors.New("password must not exceed 72 UTF-8 bytes")
	}
	return nil
}
