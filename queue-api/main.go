package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/cors"
	"net/http"
	"sync"
)


var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex
var userQueues = make(map[string]string) // sessionID -> queue

func broadcastEvent(event string) {
	mu.Lock()
	defer mu.Unlock()
	for conn := range clients {
		err := conn.WriteJSON(gin.H{"event": event})
		if err != nil {
			delete(clients, conn)
			conn.Close()
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	   r := gin.Default()
	   r.Use(cors.New(cors.Config{
		   AllowOrigins:     []string{"*"},
		   AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		   AllowHeaders:     []string{"Origin", "Content-Type", "x-client-id"},
		   ExposeHeaders:    []string{"Content-Length"},
		   AllowCredentials: true,
	   }))

	   // Helper: get client id from header (or fallback to "A0")
	   getClientID := func(c *gin.Context) string {
		   id := c.GetHeader("x-client-id")
		   if id == "" {
			   id = "A0"
		   }
		   return id
	   }

	   // HTTP: get current queue for this client id
	   r.GET("/queue", func(c *gin.Context) {
		   clientID := getClientID(c)
		   mu.Lock()
		   current, ok := userQueues[clientID]
		   if !ok {
			   // Find the highest queue in the system, then assign nextQueueNumber from that value
			   lastQueue := "A0"
			   for _, v := range userQueues {
				   if queueGreater(v, lastQueue) {
					   lastQueue = v
				   }
			   }
			   if len(userQueues) > 0 {
				   current = nextQueueNumber(lastQueue)
			   } else {
				   current = "A1"
			   }
			   userQueues[clientID] = current
		   }
		   mu.Unlock()
		   println("/queue clientID:", clientID, "queue:", current)
		   c.JSON(http.StatusOK, gin.H{"queue": current})
		   broadcastEvent("queue")
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

		   // get client id from header (or fallback to session/cookie or A0)
		   clientID := c.GetHeader("x-client-id")
		   if clientID == "" {
			   clientID, _ = c.Cookie("session_id")
			   if clientID == "" {
				   clientID = "A0"
			   }
		   }

		   for {
			   _, message, err := conn.ReadMessage()
			   if err != nil {
				   break
			   }

			   msg := string(message)

			   switch msg {
			   case "get":
					mu.Lock()

					// หา queue ล่าสุดในระบบ
					lastQueue := "A0"
					for _, v := range userQueues {
						if queueGreater(v, lastQueue) {
							lastQueue = v
						}
					}

					mu.Unlock()

					err := conn.WriteJSON(gin.H{"queue": lastQueue})
					if err != nil {
						println("write error:", err.Error())
					} else {
						println("sent latest queue:", lastQueue)
					}


			   case "next":
				   mu.Lock()
				   current, ok := userQueues[clientID]
				   if !ok {
					   current = "A0"
				   }
				   next := nextQueueNumber(current)
				   userQueues[clientID] = next
				   mu.Unlock()

				   // optionally, send new queue to this client only
				   err := conn.WriteJSON(gin.H{"queue": userQueues[clientID]})
				   if err != nil {
					   println("write error:", err.Error())
				   }

			   case "clear":
				   mu.Lock()
				   for k := range userQueues {
					   delete(userQueues, k)
				   }
				   println("userQueues after ws clear:", userQueues)
				   mu.Unlock()

				   err := conn.WriteJSON(gin.H{"queue": "A0"})
				   if err != nil {
					   println("write error:", err.Error())
				   }
			   }
		   }
	   })

	r.Run(":8080")
}

// nextQueueNumber returns next queue in A0–Z9 cycle
func nextQueueNumber(current string) string {
	if len(current) != 2 {
		return "A0"
	}
	c := current[0]
	n := current[1]

	if c == 'Z' && n == '9' {
		return "A0"
	} else if n == '9' {
		return string([]byte{c + 1, '0'})
	} else {
		return string([]byte{c, n + 1})
	}
}
// queueGreater เปรียบเทียบคิวแบบ A0–Z9 ว่า a > b
func queueGreater(a, b string) bool {
   if len(a) != 2 || len(b) != 2 {
	   return false
   }
   if a[0] > b[0] {
	   return true
   }
   if a[0] == b[0] && a[1] > b[1] {
	   return true
   }
   return false
}