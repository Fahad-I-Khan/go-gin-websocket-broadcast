package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections
	},
}

type WebSocketHub struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	mutex     sync.Mutex
}

func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
	}
}

func (hub *WebSocketHub) Run() {
	for {
		message := <-hub.broadcast
		hub.mutex.Lock()
		for client := range hub.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error sending message:", err)
				client.Close()
				delete(hub.clients, client)
			}
		}
		hub.mutex.Unlock()
	}
}

func (hub *WebSocketHub) HandleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	hub.mutex.Lock()
	hub.clients[conn] = true
	hub.mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			hub.mutex.Lock()
			delete(hub.clients, conn)
			hub.mutex.Unlock()
			break
		}
		hub.broadcast <- msg
	}
}

func main() {
	r := gin.Default()
	hub := NewWebSocketHub()
	go hub.Run()

	r.GET("/ws", hub.HandleConnections)

	r.Run(":8080")
}
