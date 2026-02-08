package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gin-contrib/cors"
	"net/http"
	"sync"
	"queue-api/db"
	"database/sql"
	"os"
)

var clients = make(map[*websocket.Conn]bool)
var mu sync.Mutex
var userQueues = make(map[string]string) // sessionID -> queue
var sqliteDB *sql.DB

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

	// Initialize SQLite DB
	var err error
	sqliteDB, err = db.InitDB("queue.db")
	if err != nil {
		panic("Failed to open SQLite DB: " + err.Error())
	}

	// Auto-create schema if not exists
	schema, err := os.ReadFile("db/schema.sql")
	if err == nil {
		_, err = sqliteDB.Exec(string(schema))
		if err != nil {
			println("Failed to initialize DB schema:", err.Error())
		}
	} else {
		println("Could not read schema.sql:", err.Error())
	}

	   // Helper: get client id from header (or fallback to "A0")
	   getClientID := func(c *gin.Context) string {
		   id := c.GetHeader("x-client-id")
		   if id == "" {
			   id = "A0"
		   }
		   return id
	   }

	   // HTTP: get current active queue for this client id
	   r.GET("/queue", func(c *gin.Context) {
		   clientID := getClientID(c)
		   // Try to get active queue from DB
		   row := sqliteDB.QueryRow("SELECT queue_number FROM queue WHERE client_id = ? AND status = 'active' ORDER BY id DESC LIMIT 1", clientID)
		   var current string
		   err := row.Scan(&current)
		   if err != nil {
			   // Assign new queue if no active
			   rows, err := sqliteDB.Query("SELECT queue_number FROM queue WHERE status = 'active' ORDER BY id DESC LIMIT 1")
			   lastQueue := "A0"
			   if err == nil && rows.Next() {
				   var last string
				   rows.Scan(&last)
				   lastQueue = last
			   }
			   rows.Close()
			   if lastQueue != "A0" {
				   current = nextQueueNumber(lastQueue)
			   } else {
				   current = "A1"
			   }
			   db.AddQueue(sqliteDB, clientID, current)
		   }
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
				// Find last active queue in DB
				row := sqliteDB.QueryRow("SELECT queue_number FROM queue WHERE status = 'active' ORDER BY id DESC LIMIT 1")
				lastQueue := "A0"
				err := row.Scan(&lastQueue)
				if err != nil {
					lastQueue = "A0"
				}
				err = conn.WriteJSON(gin.H{"queue": lastQueue})
				if err != nil {
					println("write error:", err.Error())
				}

			case "next":
				// Get current queue for client
				q, err := db.GetQueueByClientID(sqliteDB, clientID)
				current := "A0"
				if err == nil && q != nil {
					current = q.QueueNumber
				}
				next := nextQueueNumber(current)
				db.AddQueue(sqliteDB, clientID, next)
				err = conn.WriteJSON(gin.H{"queue": next})
				if err != nil {
					println("write error:", err.Error())
				}

			case "clear":
				db.ClearQueues(sqliteDB)
				err = conn.WriteJSON(gin.H{"queue": "A0"})
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