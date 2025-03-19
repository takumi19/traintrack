package main

import (
	"fmt"
	"net/http"
	"strings"
	"traintrack/internal/database"

	"github.com/golang-jwt/jwt/v5"
)

// This one checks the JWT token for the request from the header if it is present, if it is it checks it
// and passes it down through the request context.
func (a *Api) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Tells the cache to cache responses based on this header
		// NOTE: Might not be useful for a web frontend, but sending jwt in a header for a mobile app is fine
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		// If there is no token there is nothing to validate
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) == 2 && headerParts[0] == "Bearer" {
			tokenString := headerParts[1]

			token, err := validateJWT(tokenString)
			if err != nil {
				a.l.Level(DEBUG).Println(err.Error())
				WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Failed to validate JWT"})
				return
			}

			if !token.Valid {
				w.Header().Add("WWW-Authenticate", "Bearer")
				WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Failed to validate JWT"})
				return
			}

			sub, err := token.Claims.GetSubject()
			if err != nil {
				WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "internal"})
				return
			}

			u, err := a.db.GetUserByEmail(sub)
			if err != nil || u == nil {
				a.l.Level(WARN).Println("Failed to find user form the JWT email")
				WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "internal"})
				return
			}

			// set context authenticated user
			r = ctxSetAuthenticatedUser(r, u)

			// ServeHTTP
		}
		next.ServeHTTP(w, r)
	})
}

func (a *Api) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ctxGetAuthenticatedUser(r.Context())

		if user == nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "You must be authenticated to access this resource"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func createJWT(user *database.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"sub":       user.Email,
	}

	// secret := os.Getenv("JWT_SECRET")
	secret := "temporary_secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	// secret := os.Getenv("JWT_SECRET")
	secret := "temporary_secret"

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
