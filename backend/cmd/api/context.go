package main

import (
	"context"
	"net/http"
	"traintrack/internal/database"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func ctxSetAuthenticatedUser(r *http.Request, u *database.User) *http.Request {
  newCtx := context.WithValue(r.Context(), authenticatedUserContextKey, u)
  return r.WithContext(newCtx)
}

func ctxGetAuthenticatedUser(ctx context.Context) *database.User {
  u, ok := ctx.Value(authenticatedUserContextKey).(*database.User)
  if !ok {
    return nil
  }

  return u
}
