package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client should send the type of element that was changed and then the changed state
// See https://go.dev/blog/json for tips on decoding an unknown json
func (a *Api) handleEditProgram(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		_ = fmt.Errorf("websocket upgrader: %w", err)
		return
	}
	fmt.Printf("Handling the websocket connection...\n")
	defer fmt.Println("Handler returned")

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if msgType == websocket.BinaryMessage {
			return
		}
		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

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
		log.Default().Println(err)
		WriteJSON(w, http.StatusInternalServerError, &ApiError{Error: "something went wrong"})
		return
	}

	if err := WriteJSON(w, http.StatusOK, programs); err != nil {
		log.Default().Println("Failed to encode JSON:", err)
	}
}

func (a *Api) handleCreateProgram(w http.ResponseWriter, r *http.Request) {}

func (a *Api) handleReadProgram(w http.ResponseWriter, r *http.Request) {}

func (a *Api) handleUpdateProgram(w http.ResponseWriter, r *http.Request) {}

func (a *Api) handleDeleteProgram(w http.ResponseWriter, r *http.Request) {}
