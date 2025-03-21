package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"traintrack/internal/database"
	"traintrack/internal/password"
	"traintrack/internal/jwt"
)

func (a *Api) handleLogin(w http.ResponseWriter, r *http.Request) {
	// Validate the password from the data in the database
	// If not valid return unauthorized
	// If valid, creates a JWT token and sends it back

	var data userDTO
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		a.l.Level(ERROR).Print("Failed to decode JSON from the request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

  if err := data.validate(); err != nil {
    a.l.Level(ERROR).Println(err)
		WriteJSON(w, http.StatusBadRequest, ApiError{Error: "missing fields"})
		return
	}

	encryptedPassword, err := password.Hash(*data.PlaintextPassword)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "Hashing problems"})
		return
	}

	user := &database.User{
		FullName:     data.FullName,
		Login:        data.Login,
		Email:        data.Email,
		PasswordHash: &encryptedPassword,
	}

	existingUser, err := a.db.GetUserByEmail(*data.Email)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			WriteJSON(w, http.StatusNotFound, &ApiError{Error: "email was not found"})
			return
		}
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "internal"})
		return
	}

	matches, err := password.Matches(*data.PlaintextPassword, *existingUser.PasswordHash)
	if err != nil {
		a.l.Level(ERROR).Println("Failed to hash the password")
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "internal"})
		return
	}

	if !matches {
		WriteJSON(w, http.StatusUnauthorized, &ApiError{Error: "Wrong password"})
		return
	}

	// Create JWT
	token, err := jwt.CreateJWT(a.c.jwt.secretKey, user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "Failed to create token"})
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
