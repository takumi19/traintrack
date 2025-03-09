package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Connection hub for the websocket connections
var cHub map[int64]*websocket.Conn

// Client should send the type of element that was changed and then the changed state
// See https://go.dev/blog/json for tips on decoding an unknown json
// Currently expects the whole program JSON
func (a *Api) handleEditProgram(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		_ = fmt.Errorf("websocket upgrader: %w", err)
		return
	}
	a.l.Level(INFO).Println("Handling the websocket connection")
	defer a.l.Level(INFO).Println("Handler returned")

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

    // if err := json.NewDecoder(msg).Decode(program)

		a.l.Level(INFO).Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

// Returns all programs with all weeks, days, exercises and sets
func (a *Api) handleListPrograms(w http.ResponseWriter, r *http.Request) {
	programs, err := a.s.ListPrograms()
	if err != nil {
		a.l.Level(ERROR).Print(err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "something went wrong"})
		return
	}

	if err := WriteJSON(w, http.StatusOK, programs); err != nil {
		a.l.Level(ERROR).Print("Failed to encode JSON:", err)
	}
}

func (a *Api) handleCreateProgram(w http.ResponseWriter, r *http.Request) {}

func (a *Api) handleReadProgram(w http.ResponseWriter, r *http.Request) {}

func (a *Api) handleUpdateProgram(w http.ResponseWriter, r *http.Request) {}

func (a *Api) handleDeleteProgram(w http.ResponseWriter, r *http.Request) {}
