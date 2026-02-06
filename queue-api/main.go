package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

var queue []string
var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex
// Broadcast the queue to all connected clients
func broadcastQueue() {
	mu.Lock()
	defer mu.Unlock()
	for client := range clients {
		err := client.WriteJSON(gin.H{"queue": queue, "count": len(queue)})
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// main initializes the Gin web server and sets up HTTP and WebSocket endpoints for a real-time queue API.
// 
// Endpoints:
//   - POST /enqueue: Adds an item to the queue. Expects a JSON body with an "item" field.
//   - POST /dequeue: Removes and returns the first item from the queue. If the queue is empty, returns nil.
//   - GET /queue: Returns the current state of the queue.
//   - GET /ws: Upgrades the connection to a WebSocket for real-time queue updates. Clients can send messages to add items or send "clear" to reset the queue. All connected clients receive updates via broadcast.
//
// The function manages concurrent access to the queue and connected WebSocket clients using a mutex.
func main() {
	r := gin.Default()

	// Add item to queue
	r.POST("/enqueue", func(c *gin.Context) {
		var json struct {
			Item string `json:"item" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		queue = append(queue, json.Item)
		c.JSON(http.StatusOK, gin.H{"queue": queue})
	})

	// Remove item from queue
	r.POST("/dequeue", func(c *gin.Context) {
		if len(queue) == 0 {
			c.JSON(http.StatusOK, gin.H{"item": nil, "queue": queue})
			return
		}
		item := queue[0]
		queue = queue[1:]
		c.JSON(http.StatusOK, gin.H{"item": item, "queue": queue})
	})

	// Get current queue
	r.GET("/queue", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"queue": queue})
	})

	// WebSocket endpoint
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		mu.Lock()
		clients[conn] = true
		mu.Unlock()
		defer func() {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close()
		}()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}
			msg := string(message)
			mu.Lock()
			if msg == "clear" {
				queue = []string{}
			} else {
				queue = append(queue, msg)
			}
			mu.Unlock()
			
			broadcastQueue()
		}
	})

	r.Run(":8080")
}
