package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"traintrack/internal/database"
)

func (a *Api) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.db.ListUsers()
	if err != nil {
		a.l.Level(ERROR).Print("Failed to list all users:", err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "Database issues"})
		return
	}
	WriteJSON(w, http.StatusOK, users)
}

// BUG: This and several others are not idempotent, use Idempotency-key header in the request to handle this
func (a *Api) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-type"); contentType != "application/json" {
		a.l.Level(ERROR).Print("Wrong Content-type:", contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		a.l.Level(ERROR).Print("Empty request body when creating user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user database.User
	// TODO: There is nuance to handling decoder errors
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.l.Level(ERROR).Print("Failed to decode JSON from the request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := a.db.CreateUser(&user)
	if err != nil {
		errMsg := "Failed to create new user"
		a.l.Level(ERROR).Print(errMsg, err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: errMsg})
		return
	}

	WriteJSON(w, http.StatusOK, struct{ Id int64 }{Id: id})
}

func (a *Api) handleReadUser(w http.ResponseWriter, r *http.Request) {
	// NOTE: id should be the last part of the url
	paths := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(paths[len(paths)-1])
	if err != nil {
		a.l.Level(ERROR).Println("Invalid id:", paths[len(paths)-1])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := a.db.ReadUser(int64(id))
	if err != nil {
		a.l.Level(ERROR).Println("Failed to read user data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = WriteJSON(w, http.StatusFound, user); err != nil {
		a.l.Level(ERROR).Println("Failed to convert user to JSON:", err)
	}
}

func (a *Api) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the id - should be the last part of the url
	paths := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(paths[len(paths)-1])
	if err != nil {
		a.l.Level(ERROR).Println("Invalid id:", paths[len(paths)-1])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = a.db.DeleteUser(int64(id)); err != nil {
		a.l.Level(ERROR).Println("Failed to delete user with id", id)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "db fault"})
	}
}

func (a *Api) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	a.l.Level(INFO).Println("Updating user info")

	if contentType := r.Header.Get("Content-type"); contentType != "application/json" || r.Body == nil {
		a.l.Level(ERROR).Println("Bad request when updating user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the id - should be the last part of the url
	paths := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(paths[len(paths)-1])
	if err != nil {
		a.l.Level(ERROR).Println("Invalid id:", paths[len(paths)-1])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Decode JSON into the struct
	var user database.User
	// TODO: There is nuance to handling decoder errors
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		a.l.Level(ERROR).Println("Failed to decode JSON from the request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Put this here bc the decoder may overwrite the id.
	// But i need to double check that.
	user.Id = id
	if err = a.db.UpdateUser(&user); err != nil {
		a.l.Level(ERROR).Println(err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "db fault"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
