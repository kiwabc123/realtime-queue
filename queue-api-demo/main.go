package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	queue   []string
	mu      sync.Mutex
	clients = make(map[*websocket.Conn]bool)
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// ส่งคิวไปให้ client ทุกตัว
func broadcastQueue() {
	mu.Lock()
	defer mu.Unlock()

	for client := range clients {
		err := client.WriteJSON(map[string]interface{}{"queue": queue ,"count": len(queue)})
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	clients[conn] = true

	// ส่งคิวตอนเชื่อมต่อ
	conn.WriteJSON(map[string]interface{}{"queue": queue ,"count": len(queue)})

	for {
		var msg map[string]string
		err := conn.ReadJSON(&msg)
		if err != nil {
			delete(clients, conn)
			conn.Close()
			break
		}

		mu.Lock()
		switch msg["action"] {
		case "add":
			queue = append(queue, msg["value"])
		case "clear":
			queue = []string{}
		}
		mu.Unlock()

		broadcastQueue()
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/ws", wsHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}


//docker run --rm -it -v ${PWD}:/app -w /app -p 8080:8080 golang:1.21 sh -c "go get github.com/gorilla/websocket && go run main.go"