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

func (a *Api) debugHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is a protected handler"))
}

func (api *Api) trainingLogRoutes() http.Handler {
  mux := http.NewServeMux()
	// mux.HandleFunc("GET /logs", api.handleListLogs)
	// mux.HandleFunc("GET /logs/{log_id}", api.handleReadLog)
	// mux.HandleFunc("PATCH /logs/{log_id}", api.handleUpdateLog)
	// mux.HandleFunc("DELETE /logs/{log_id}", api.handleDeleteProgram)
	// mux.HandleFunc("GET /logs/{log_id}/{workout_id}", api.handleReadLoggedWorkout)
	// mux.HandleFunc("POST /logs/{log_id}/{workout_id}", api.handleCreateLoggedWorkout)
  return mux
}

func (api *Api) userRoutes() http.Handler {
  mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.handleListUsers)
	mux.HandleFunc("POST /", api.handleCreateUser)
	mux.HandleFunc("GET /{user_id}", api.handleReadUser)
	mux.HandleFunc("DELETE /{user_id}", api.handleDeleteUser)
	mux.HandleFunc("PATCH /{user_id}", api.handleUpdateUser)
  return mux
}

func (api *Api) chatRoutes() http.Handler {
  mux := http.NewServeMux()
	mux.HandleFunc("/", api.handleListUserChats)
	mux.HandleFunc("/{chat_id}", api.handleChatWs)
  return mux
}

func (api *Api) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /test", api.requireAuthentication(http.HandlerFunc(api.debugHandle)))

	mux.HandleFunc("POST /login", api.handleLogin)
	mux.HandleFunc("POST /signup", api.handleSignup)

  mux.Handle("/users/", http.StripPrefix("/users", api.userRoutes()))

	mux.HandleFunc("GET /programs", api.handleListPrograms)
	mux.HandleFunc("GET /programs/{template_id}/edit", api.handleEditProgram)
	// mux.Handle("GET /programs/{template_id}/edit", api.requireAuthentication(http.HandlerFunc(api.handleEditProgram)))

  mux.Handle("/logs/", http.StripPrefix("/logs", api.trainingLogRoutes()))

  // mux.Handle("/chats/", http.StripPrefix("/chats", api.chatRoutes()))
	mux.HandleFunc("/chats", api.handleListUserChats)
	mux.HandleFunc("/chats/{chat_id}", api.handleChatWs)

	mux.HandleFunc("GET /exercises", api.handleListExercises)
	mux.HandleFunc("GET /exercises/{exercise_id}", api.handleListExercises)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	return v1
}
