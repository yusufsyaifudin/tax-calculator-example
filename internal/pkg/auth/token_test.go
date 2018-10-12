package auth

import (
	"context"
	"fmt"
	"testing"
)

const secretKey = "mysecret"
const userID = 1

func TestGenerateJWTToken(t *testing.T) {
	_, err := GenerateJWTToken(context.Background(), secretKey, userID)
	if err != nil {
		t.Error(err)
	}
}

func TestValidateJWTToken(t *testing.T) {
	ctx := context.Background()
	token, err := GenerateJWTToken(ctx, secretKey, userID)
	if err != nil {
		t.Error(err)
	}

	userIDFromToken, err := ValidateJWTToken(ctx, secretKey, token)
	if err != nil {
		t.Error(err)
	}

	if userID != userIDFromToken {
		t.Error(fmt.Errorf("error user id in token != generated"))
	}
}
