package main

import (
	"encoding/json"
	"net/http"
	"traintrack/internal/chat"
	"traintrack/internal/database"
	"traintrack/internal/editor"
	"traintrack/internal/middleware"
)

type Api struct {
	db *database.DB
	l  Logger
	c  config
	// Editor hub
	eHub *editor.Hub
	// Chat hub
	cHub *chat.Hub
}

type ApiError struct {
	Error string
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

func (api *Api) exerciseRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.handleListExercises)
	mux.HandleFunc("GET /{exercise_id}", api.handleGetExerciseByID)
	return mux
}

func (api *Api) chatRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", api.handleListUserChats)
	mux.HandleFunc("/{chat_id}", api.handleChatWs)
	return mux
}

func (api *Api) programRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", api.handleListPrograms)
	mux.HandleFunc("GET /{template_id}/edit", api.handleEditProgram)
	return mux
}

func (api *Api) routes() http.Handler {
	mux := http.NewServeMux()

	// Testing route
	mux.Handle("POST /test", middleware.RequireAuthentication(http.HandlerFunc(api.debugHandle)))

	// Login and signup handlers
	mux.HandleFunc("POST /login", api.handleLogin)
	mux.HandleFunc("POST /signup", api.handleSignup)

	// User handlers
	mux.Handle("/users/", http.StripPrefix("/users", api.userRoutes()))

	// Program templates handlers
	mux.Handle("/programs/", http.StripPrefix("/programs", api.programRoutes())) // TODO: Wrap these routes in a check program rights auth middleware

	// Program logs handlers
	mux.Handle("/logs/", http.StripPrefix("/logs", api.trainingLogRoutes()))

	// Chat handlers
	mux.Handle("/chats/", http.StripPrefix("/chats", api.chatRoutes()))

	// Exercise info handlers
	mux.Handle("/exercises/", http.StripPrefix("/exercises", api.exerciseRoutes()))

	// Separate API versions
	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", mux))

	return v1
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
