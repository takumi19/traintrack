package jwt

import (
	"fmt"
	"time"
	"traintrack/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

// var ErrTokenExpired = jwt.ErrTokenExpired

type CustomClaims struct {
	jwt.RegisteredClaims
}

func NewAccessToken(secret string, user *database.User) (string, error) {
	claims := CustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   *user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

// Contains:
//
//	Version - kept in the DB, incremented upon logout
//	User ID
//
// Is signed using the user's password so that if the user resets their password,
// the old tokens are no longer valid
func NewRefreshToken(secret string, user *database.User) (string, error) {
	claims := CustomClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   *user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func GetToken(secret string, tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
}

// func CheckAccessToken(accessSecret, accessToken, refreshSecret, refreshToken string) error {
//   token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
//
// 		return []byte(accessSecret), nil
// 	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
//
//   // If the token is valid we just return
//   if token.Valid {
//     return nil
//   }
//
//   if errors.Is(err, jwt.ErrTokenExpired) {
//     return err
//   }
//
//   return nil
// }
