package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"traintrack/internal/database"
	"traintrack/internal/jwt"
)

type AuthDeps struct {
	DB        *database.DB
	Logger    *slog.Logger
	JwtSecret string
}

// AuthMiddleware returns an http.Handler that requires authentication.
func AuthMiddleware(deps AuthDeps) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

				token, err := jwt.GetToken(deps.JwtSecret, tokenString)
				if err != nil {
					authFailure(w, "Failed to validate JWT")
					return
				}

				if !token.Valid {
					authFailure(w, "Failed to validate JWT")
					return
				}

				sub, err := token.Claims.GetSubject()
				if err != nil {
					authFailure(w, "Failed to parse JWT")
					return
				}

				u, err := deps.DB.GetUserByEmail(sub)
				if err != nil || u == nil {
					deps.Logger.Warn("Failed to find user form the JWT email")
					authFailure(w, "Failed to find user by email")
					return
				}

				// set context authenticated user
				r = ctxSetAuthenticatedUser(r, u)
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := ctxGetAuthenticatedUser(r)

		if user == nil {
      authFailure(w, "Failed to authenticate user")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authFailure(w http.ResponseWriter, msg string) {
	w.Header().Add("Content-Type", "application/json")
  w.Header().Add("WWW-Authenticate", "Bearer")
  w.WriteHeader(http.StatusUnauthorized)

	data := struct {
		Error string
	}{Error: msg}

	json.NewEncoder(w).Encode(data)
}
