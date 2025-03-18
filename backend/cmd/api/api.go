package main

import (
	"encoding/json"
	"net/http"
	"traintrack/internal/chat"
	"traintrack/internal/database"
	"traintrack/internal/editor"
)

type Api struct {
	db *database.DB
	l  Logger
	// Editor hub
	eHub *editor.Hub
	// Chat hub
	cHub *chat.Hub
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

// TODO: Rewrite this mf
func (api *Api) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /test", api.testhndlr)

	mux.HandleFunc("GET /users", api.handleListUsers)
	mux.HandleFunc("POST /users", api.handleCreateUser)
	// NOTE: Request.PathValue matches the {user_id}
	mux.HandleFunc("GET /users/{user_id}", api.handleReadUser)
	mux.HandleFunc("DELETE /users/{user_id}", api.handleDeleteUser)
	mux.HandleFunc("PATCH /users/{user_id}", api.handleUpdateUser)

	mux.HandleFunc("GET /programs", api.handleListPrograms)
	mux.HandleFunc("GET /programs/{template_id}/edit", api.handleEditProgram)
	// mux.HandleFunc("GET /programs/{program_id}/edit", handleEditProgram)

	// mux.HandleFunc("GET /logs", api.handleListLogs)
	// mux.HandleFunc("GET /logs/{log_id}", api.handleReadLog)
	// mux.HandleFunc("PATCH /logs/{log_id}", api.handleUpdateLog)
	// mux.HandleFunc("DELETE /logs/{log_id}", api.handleDeleteProgram)
	// mux.HandleFunc("GET /logs/{log_id}/{workout_id}", api.handleReadLoggedWorkout)
	// mux.HandleFunc("POST /logs/{log_id}/{workout_id}", api.handleCreateLoggedWorkout)

	// mux.HandleFunc("POST /chats", api.handleCreateChat)
  // mux.HandleFunc("GET /chats", api.handleListChats)
	// mux.HandleFunc("GET /chats/{chat_id}", api.handleReadChat)
	// mux.HandleFunc("DELETE /chats/{chat_id}", api.handleDeleteChat)

  // For the chats we need:
  // GET chats by user_id
  // DELETE chat by chat_id
  // Websocket connection to concrete chat
	mux.HandleFunc("/chats/{chat_id}", api.handleChatWs)
	// mux.HandleFunc("DELETE /chats/{chat_id}", api.handleDeleteChatMessage)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	return v1
}
