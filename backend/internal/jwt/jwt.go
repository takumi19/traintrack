package jwt

import (
	"fmt"
	"traintrack/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret string, user *database.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"sub":       user.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func ValidateJWT(secret string, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
