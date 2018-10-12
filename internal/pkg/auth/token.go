package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gbrlsnchs/jwt"
)

const issuer = "tax-calculator-example"
const audience = "user"

// GenerateJWTToken will generate JWT token based on input user.
// Returns token, error.
func GenerateJWTToken(parent context.Context, secretKey string, userId int64) (string, error) {
	_, cancel := context.WithTimeout(parent, time.Duration(1)*time.Second)
	defer cancel()

	// Timestamp the beginning.
	now := time.Now()

	userIdStr := fmt.Sprintf("%d", userId)

	// Define a signer.
	hs256 := jwt.NewHS256(secretKey)
	jot := &jwt.JWT{
		Issuer:         issuer,
		Subject:        userIdStr,
		Audience:       audience,
		ExpirationTime: now.Add(24 * 30 * 12 * time.Hour).Unix(), // token is valid for 1 year
		NotBefore:      now.Unix(),                               // token can be used right now once it generated
		IssuedAt:       now.Unix(),
		ID:             userIdStr,
	}

	jot.SetAlgorithm(hs256)
	payload, err := jwt.Marshal(jot)
	if err != nil {
		return "", err
	}

	tokenBytes, err := hs256.Sign(payload)
	if err != nil {
		return "", err
	}

	return string(tokenBytes), nil
}

// ValidateJWTToken will return user model if it success.
func ValidateJWTToken(parent context.Context, secretKey, token string) (userId int64, err error) {
	_, cancel := context.WithTimeout(parent, time.Duration(1)*time.Second)
	defer cancel()

	now := time.Now()
	hs256 := jwt.NewHS256(secretKey) // Define a signer.

	// First, extract the payload and signature.
	// This enables un-marshaling the JWT first and verifying it later or vice versa.
	payload, sig, err := jwt.Parse(token)
	if err != nil {
		return 0, err
	}

	if err = hs256.Verify(payload, sig); err != nil {
		return 0, err
	}

	var jot jwt.JWT
	if err = jwt.Unmarshal(payload, &jot); err != nil {
		return 0, err
	}

	// Validate fields.
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now)
	audValidator := jwt.AudienceValidator(audience)
	issValidator := jwt.IssuerValidator(issuer)
	err = jot.Validate(iatValidator, expValidator, audValidator, issValidator)
	if err != nil {
		return 0, err
	}

	userIdInt, err := strconv.Atoi(jot.Subject)
	if err != nil {
		return 0, err
	}

	return int64(userIdInt), nil
}
