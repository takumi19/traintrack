package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"traintrack/internal/database"
	"traintrack/internal/password"
	"traintrack/internal/jwt"
)

type userDTO struct {
	FullName          *string `json:"full_name" db:"full_name"`
	Login             *string `json:"login" db:"login"`
	Email             *string `json:"email" db:"email"`
	PlaintextPassword *string `json:"password" db:"password"`
}

func (a *Api) handleSignup(w http.ResponseWriter, r *http.Request) {
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

	_, err = a.db.CreateUser(user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "Failed to create user"})
		return
	}

	token, err := jwt.CreateJWT(a.c.jwt.secretKey, user)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, ApiError{Error: "Failed to create token"})
		return
	}

	WriteJSON(w, http.StatusCreated, map[string]string{
		"token": token,
	})
}

func (data *userDTO) validate() error {
	if data.Email == nil {
		return errors.New("no email")
	}
	if data.FullName == nil {
		return errors.New("no full name")
	}
	if data.Login == nil {
		return errors.New("no login")
	}
	if data.PlaintextPassword == nil {
		return errors.New("no password")
	}
	return nil
}
