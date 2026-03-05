package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Secret = []byte("DocNebula-secret")

// Login JWT
func GenerateToken(userID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(Secret)
}

// Reset password token
func GenerateResetToken(userID string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(30 * time.Minute).Unix(),
		"type":    "reset",
	})

	return token.SignedString(Secret)
}

// Verify reset token
func VerifyResetToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
	})

	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	userID := claims["user_id"].(string)

	return userID, nil
}
