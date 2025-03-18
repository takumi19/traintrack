package main

import (
	"net/http"
	"traintrack/internal/chat"
	// "github.com/gorilla/websocket"
)

// var upgrader = websocket.Upgrader{
//   ReadBufferSize: 1024,
//   WriteBufferSize: 1024,
// }

const (
	messageSize = 512
)

func (a *Api) handleListUserChats(w http.ResponseWriter, r *http.Request) {
	// chatId, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	// if err != nil {
	// 	WriteJSON(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }
}

func (a *Api) handleChatWs(w http.ResponseWriter, r *http.Request) {
	// chatId, err := strconv.ParseInt(r.URL.Query().Get("user_id"), 10, 64)
	// if err != nil {
	// 	WriteJSON(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		a.l.Level(FATAL).Fatal("Failed to establish a websocket connection for the chats")
		return
	}

	client := &chat.Client{Hub: a.cHub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
