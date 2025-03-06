package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (a *Api) handleListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.s.ListUsers()
	if err != nil {
		log.Default().Println("Failed to list all users:", err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "Database issues"})
		return
	}
	WriteJSON(w, http.StatusOK, users)
}

// BUG: This and several others are not idempotent, use Idempotency-key header in the request to handle this
func (a *Api) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if contentType := r.Header.Get("Content-type"); contentType != "application/json" {
		log.Default().Println("Wrong Content-type:", contentType)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Body == nil {
		log.Default().Println("Empty request body when creating user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	// TODO: There is nuance to handling decoder errors
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Default().Println("Failed to decode JSON from the request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := a.s.CreateUser(&user)
	if err != nil {
		errMsg := "Failed to create new user"
		log.Default().Println(errMsg, err)
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
		log.Default().Println("Invalid id:", paths[len(paths)-1])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := a.s.ReadUser(int64(id))
	if err != nil {
		log.Default().Println("Failed to read user data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = WriteJSON(w, http.StatusFound, user); err != nil {
		log.Default().Println("Failed to convert user to JSON:", err)
	}
}

func (a *Api) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Get the id - should be the last part of the url
	paths := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(paths[len(paths)-1])
	if err != nil {
		log.Default().Println("Invalid id:", paths[len(paths)-1])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = a.s.DeleteUser(int64(id)); err != nil {
		log.Default().Println("Failed to delete user with id", id)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "db fault"})
	}
}

func (a *Api) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("Updating user info")

	if contentType := r.Header.Get("Content-type"); contentType != "application/json" || r.Body == nil {
		log.Default().Println("Bad request when updating user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the id - should be the last part of the url
	paths := strings.Split(r.URL.Path, "/")
	id, err := strconv.Atoi(paths[len(paths)-1])
	if err != nil {
		log.Default().Println("Invalid id:", paths[len(paths)-1])
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Decode JSON into the struct
	var user User
	// TODO: There is nuance to handling decoder errors
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Default().Println("Failed to decode JSON from the request:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: Put this here bc the decoder my overwrite the id.
	// But i need to double check that.
	user.Id = id
	if err = a.s.UpdateUser(&user); err != nil {
		log.Default().Println(err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "db fault"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
