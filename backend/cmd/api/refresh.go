package main

import (
	"errors"
	"net/http"
	"traintrack/internal/jwt"

	"github.com/gorilla/securecookie"
)

func (a *Api) handleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(a.tokenCookieTemplate.Name)
	if err != nil {
		// Only expected error is http.ErrNoCookie.
		if !errors.Is(err, http.ErrNoCookie) {
			a.l.Level(WARN).Println("unexpected error getting session cookie")
		}
    WriteError(w, http.StatusBadRequest, "Bad cookie")
		return
	}

	var tokenStr string
	err = a.secureCookie.Decode(a.tokenCookieTemplate.Name, cookie.Value, &tokenStr)
	if err != nil {
		var secureCookieError securecookie.Error
		if errors.As(err, &secureCookieError) && secureCookieError.IsDecode() {
			a.l.Level(ERROR).Println("Failed to decode cookie:", err.Error())
		}
    WriteError(w, http.StatusBadRequest, "Bad cookie")
		return
	}

	token, err := jwt.GetToken(a.c.jwt.refreshKey, tokenStr)
	if !token.Valid {
		WriteError(w, http.StatusUnauthorized, "Expired or invalid token")
		a.l.Level(INFO).Println(err.Error())
		return
	}

  email, err := token.Claims.GetSubject()
  if err != nil {
    a.l.Level(ERROR).Println(err.Error())
    WriteError(w, http.StatusInternalServerError, "Failed to parse subject of the token")
    return
  }

  user, err := a.db.GetUserByEmail(email)
  if err != nil {
    a.l.Level(ERROR).Println(err.Error())
    WriteError(w, http.StatusInternalServerError, "Invalid user id encoded in the refresh token")
    return
  }

	newAccessToken, err := jwt.NewAccessToken(a.c.jwt.accessKey, user)
  if err != nil {
    WriteError(w, http.StatusInternalServerError, "Failed to create new access token")
    a.l.Level(ERROR).Println(err.Error())
    return
  }

  WriteJSON(w, http.StatusOK, map[string]string{
    "token": newAccessToken,
  })
}
