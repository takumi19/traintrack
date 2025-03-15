package editor

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"traintrack/internal/database"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 16384
)

const (
	program         = "program_template"
	programWeek     = "program_week"
	programWorkout  = "program_workout"
	workoutExercise = "program_exercise"
	workoutSet      = "program_set"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// This is a wrapper type for the messages received from the frontend
// The type specifies whether it is a program, a week, etc. Must be
// specifiede by the client before sending.
// The Data field is the actual payload received.
type MessageWrapper struct {
	// The Id is the id of the program being edited. Set before being sent to the hub
	Id int64 `json:"-"`
	// The type of element edited as set by the frontend
	Type string `json:"type"`
	// The actual paylod received from the client
	Data json.RawMessage `json:"data"`
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan MessageWrapper

	// Program ID
	ProrgamID int64
}

func (c *Client) ReadPump(db *database.DB) {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		// Read message and return if a close error is encountered
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var convertedMsg MessageWrapper
		err = json.Unmarshal(message, &convertedMsg)
		if err != nil {
			log.Default().Println("Failed to decode the message received through websockets:", err)
		}
		// Specify the id of the program being edited
		convertedMsg.Id = c.ProrgamID

		if err := c.processMessage(&convertedMsg, db); err != nil {
			log.Default().Printf("error: %v\n", err)
			break
		}

		// not sure if i need to do this:
		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Send the program to the hub
		c.Hub.broadcast <- convertedMsg
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message.Data)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				msg := <-c.Send
				w.Write(msg.Data)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Parses the message from the client and updates the database entry
func (c *Client) processMessage(msg *MessageWrapper, db *database.DB) error {
	switch msg.Type {
	case program:
		var p *database.Program
		if err := json.Unmarshal(msg.Data, p); err != nil {
			return err
		}
		db.UpdateProgram(p)
	case programWeek:
		var p *database.ProgramWeek
		if err := json.Unmarshal(msg.Data, p); err != nil {
			return err
		}
		db.UpdateProgramWeek(p)
	case programWorkout:
		var p *database.ProgramWorkout
		if err := json.Unmarshal(msg.Data, p); err != nil {
			return err
		}
		db.UpdateProgramWorkout(p)
	case workoutExercise:
		var p *database.WorkoutExercise
		if err := json.Unmarshal(msg.Data, p); err != nil {
			return err
		}
		db.UpdateWorkoutExercise(p)
	case workoutSet:
		var p *database.WorkoutSet
		if err := json.Unmarshal(msg.Data, p); err != nil {
			return err
		}
		db.UpdateWorkoutSet(p)
	default:
		return errors.New("malformed JSON")
	}

	return nil
}
