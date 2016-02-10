package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type connection struct {
	ws     *websocket.Conn
	send   chan []byte
	client clientInfo
}

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 512

	Payload = "Hello World!"
)

var StartNonce uint64

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer c.ws.Close()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		var decoded = make(map[string]string)

		json.Unmarshal(message, &decoded)

		workTmpl, _ := template.ParseFiles("templates/generate.js")

		buf := &bytes.Buffer{}

		workTmpl.Execute(buf, struct {
			StartNonce uint64
			EndNonce   uint64
			Payload    string
		}{StartNonce, StartNonce + 50000, Payload})

		work := string(buf.Bytes())

		StartNonce += 50000

		switch decoded["type"] {

		case "need":
			msg, _ := json.Marshal(struct {
				Type string `json:"type"`
				JS   string `json:"js"`
			}{"work", work})

			c.send <- msg

		case "payload":
			log.Println("Payload Message: ", decoded["payload"])

		default:
			log.Println("unknown message received: ", decoded)
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving WebSocket to ", r.RemoteAddr)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Write404(w, r)
		log.Println("error serving websocket: ", err)
		return
	}

	c := &connection{send: make(chan []byte, 256), ws: ws, client: clientList[r.RemoteAddr]}

	go c.writePump()
	c.readPump()
}
