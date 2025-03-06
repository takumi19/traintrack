package main

import (
	"encoding/json"
	"net/http"
)

type Api struct {
	s Storage
	l       Logger
}

type ApiError struct {
	Error string
}

func (a *Api) fail(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")

	data := struct {
		Error string
	}{Error: msg}

	resp, _ := json.Marshal(data)
	w.WriteHeader(status)
	w.Write(resp)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func (a *Api) ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		a.fail(w, "oops something evil has happened", 500)
		return
	}
	w.Write(resp)
}

func setUpRoutes(api *Api) {
	// http.HandleFunc("GET /test_websockets", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./websockets.html")
	// })

	// http.HandleFunc("/echo", handleWebsocketConnection)

	// http.HandleFunc("GET /param", handleNonParametrized)
	// http.HandleFunc("GET /param/{id}", handleParametrized)
	//
	// http.HandleFunc("GET /login", handleLogin)
	// http.HandleFunc("GET /signup", handleSignup)

	http.HandleFunc("GET /v1/users", api.handleListUsers)
	http.HandleFunc("POST /v1/users", api.handleCreateUser)
	http.HandleFunc("GET /v1/users/{user_id}", api.handleReadUser)
	http.HandleFunc("DELETE /v1/users/{user_id}", api.handleDeleteUser)
	http.HandleFunc("PATCH /v1/users/{user_id}", api.handleUpdateUser)

	http.HandleFunc("GET /v1/programs", api.handleListPrograms)
	// http.HandleFunc("POST /v1/programs", handleCreateProgram)
	http.HandleFunc("/v1/programs/{program_id}", api.handleEditProgram)
	// http.HandleFunc("PATCH /v1/programs/{program_id}", handleUpdateProgram)
	// http.HandleFunc("DELETE /v1/programs/{program_id}", handleDeleteProgram)

	// http.HandleFunc("GET /v1/programs/{program_id}/edit", handleEditProgram)
	//
	// http.HandleFunc("GET /v1/logs", handleListLogs)
	// http.HandleFunc("GET /v1/logs/{log_id}", handleReadLog)
	// http.HandleFunc("PATCH /v1/logs/{log_id}", handleUpdateLog)
	// http.HandleFunc("DELETE /v1/logs/{log_id}", handleDeleteProgram)
	// http.HandleFunc("GET /v1/logs/{log_id}/{workout_id}", handleReadLoggedWorkout)
	// http.HandleFunc("POST /v1/logs/{log_id}/{workout_id}", handleCreateLoggedWorkout)
	//
	// http.HandleFunc("GET /v1/chats", handleListChats)
	// http.HandleFunc("POST /v1/chats", handleCreateChat)
	// http.HandleFunc("GET /v1/chats/{chat_id}", handleReadChat)
	// http.HandleFunc("DELETE /v1/chats/{chat_id}", handleDeleteChat)
	//
	// http.HandleFunc("GET /v1/chats/{chat_id}/messages", handleListChatMessages)
	// http.HandleFunc("POST /v1/chats/{chat_id}/messages", handleCreateChatMessage)
	// http.HandleFunc("GET /v1/chats/{chat_id}/messages/{message_id}", handleReadChatMessage)
	// http.HandleFunc("PATCH /v1/chats/{chat_id}/messages/{message_id}", handleUpdateChatMessage)
	// http.HandleFunc("DELETE /v1/chats/{chat_id}/messages/{message_id}", handleDeleteChatMessage)
}
