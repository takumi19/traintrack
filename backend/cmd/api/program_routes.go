package main

import (
	"fmt"
	"net/http"
	"strconv"
	"traintrack/internal/editor"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// TODO:
	// WriteBufferPool: ...
	// Error: ...
}

// Currently expects the whole program JSON
func (a *Api) handleEditProgram(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Query().Get("template_id")) == 0 {
		WriteJSON(w, http.StatusBadRequest, &ApiError{Error: "no template id"})
    return
	}

	templateID, err := strconv.ParseInt(r.URL.Query().Get("template_id"), 10, 64)
  if err != nil {
    WriteJSON(w, http.StatusBadRequest, &ApiError{Error: "bad template id"})
    return
  }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		_ = fmt.Errorf("websocket upgrader: %w", err)
		return
	}
	a.l.Level(INFO).Println("New websocket connection")

	// Create new hub for the specified program ID if it does not exist yet

	client := &editor.Client{Hub: a.eHub, Conn: conn, Send: make(chan editor.MessageWrapper), ProrgamID: templateID}
	client.Hub.Register <- client

	go client.ReadPump(a.db)
	go client.WritePump()
}

// Returns all programs with all weeks, days, exercises and sets
func (a *Api) handleListPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := a.db.ListPrograms()
	if err != nil {
		a.l.Level(ERROR).Print(err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "something went wrong"})
		return
	}

	if err := WriteJSON(w, http.StatusOK, programs); err != nil {
		a.l.Level(ERROR).Print("Failed to encode JSON:", err)
	}
}
