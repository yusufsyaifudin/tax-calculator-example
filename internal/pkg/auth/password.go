package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword will hash the user password using bcrypt algorithm with cost 10
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash will check two different hash is same or not
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
