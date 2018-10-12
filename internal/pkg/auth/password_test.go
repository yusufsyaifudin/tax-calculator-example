package auth

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	_, err := HashPassword("my-password")
	if err != nil {
		t.Error(err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	var password = "my-password"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	ok := CheckPasswordHash(password, passwordHash)
	if !ok {
		t.Error(fmt.Errorf("error password mismatch"))
	}
}
