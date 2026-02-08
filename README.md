# Important: Docker Compose Version

All docker-compose commands in this guide require Docker Compose v1 (docker-compose with dash). If you are using Docker Compose v2 (docker compose with space), please adapt the commands accordingly. Official support and testing is for v1 only.


# Realtime Queue System

A simple realtime queue system using Go (Gin, WebSocket) for the backend and Vue 3 (Vite) for the frontend.

## Features
- REST API and WebSocket for queue management
- Per-client queue assignment (sequential, unique)
- Real-time queue updates and clearing via WebSocket
- Dockerized for easy deployment
- Modern, responsive frontend UI

## Project Structure
```
queue-api/         # Go backend (Gin, WebSocket)
queue-frontend/    # Vue 3 frontend (Vite)
queue-api-demo/    # (Optional) Demo or test code
```

## Requirements

### Backend
- Go 1.20+ (https://golang.org/dl/)
- Docker & Docker Compose (for containerized run)

### Frontend
- Node.js 18+ (https://nodejs.org/)
- npm 9+ (comes with Node.js)

## Quick Start (Docker Compose)
1. Build and run all services:
   ```sh
   docker-compose up --build
   ```
2. Frontend: http://localhost:5173
3. Backend API: http://localhost:8080

## Development

### Backend (Go)
- Location: `queue-api/`
- Run locally:
  ```sh
  cd queue-api
  go run main.go
  ```
- Endpoints:
  - `GET /queue` - Get or assign queue for client
  - `POST /queue/next` - Next queue for client
  - `POST /queue/clear` - Clear all queues (broadcast)
  - `GET /ws` - WebSocket endpoint

### Frontend (Vue 3)
- Location: `queue-frontend/`
- Run locally:
  ```sh
  cd queue-frontend
  npm install
  npm run dev
  ```

## SQLite DB Usage (queue.db)

The backend uses SQLite for persistent queue storage. Here are some common commands:

- **Inspect queue data:**
  ```sh
  sqlite3 queue-api/db/queue.db "SELECT * FROM queue;"
  ```
- **Clear all queue data:**
  ```sh
  sqlite3 queue-api/db/queue.db "DELETE FROM queue;"
  ```
- **Reset DB (delete file and re-init):**
  ```sh
  rm queue-api/db/queue.db
  # Then restart backend to auto-create schema
  ```

### Using SQLite in Docker

- Run SQL directly from host:
  ```sh
  docker-compose exec queue-api sqlite3 /app/queue.db "SELECT * FROM queue;"
  docker-compose exec queue-api sqlite3 /app/queue.db "DELETE FROM queue;"
  ```
- Enter the container shell for advanced usage:
  ```sh
  docker-compose exec queue-api sh
  # Then inside the container, run:
  sqlite3 /app/queue.db
  # Now you can use the sqlite3 CLI interactively, e.g.:
  # SELECT * FROM queue;
  # DELETE FROM queue;
  ```
## Useful docker-compose exec commands

- Enter backend container shell:
  ```sh
  docker-compose exec queue-api sh
  ```
- Show all queue data:
  ```sh
  docker-compose exec queue-api sqlite3 /app/queue.db "SELECT * FROM queue;"
  ```
- Delete all queue data:
  ```sh
  docker-compose exec queue-api sqlite3 /app/queue.db "DELETE FROM queue;"
  ```


## Quick Test with curl

You can test the backend API quickly using curl (replace <client-id> as needed):

## Bash (Linux/macOS/WSL)
Run this loop in a Bash shell:

```sh
for i in {1..5}; do
  curl http://localhost:8080/queue -H "x-client-id: $i"
  echo
done
```

## PowerShell (Windows)
Run this loop in Windows PowerShell:

```powershell
for ($i = 1; $i -le 5; $i++) {
  curl http://localhost:8080/queue -Headers @{"x-client-id"="$i"}
  Write-Host ""
}
```

```sh
# Get or assign a queue (single request)
```sh
curl http://localhost:8080/queue -H "x-client-id: 1"
```

# Loop: Test multiple clients (IDs 1 to 5)
```sh
for i in {1..5}; do
    curl http://localhost:8080/queue -H "x-client-id: $i"
    echo
done
```

You can also use browser dev tools to test WebSocket:

```js
let ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = e => console.log(e.data);
ws.onopen = () => ws.send('get');
// ws.send('next')
// ws.send('clear')
```
# Requirements

## Backend
- Go 1.20+ (https://golang.org/dl/)
- Docker & Docker Compose (for containerized run)

## Frontend
- Node.js 18+ (https://nodejs.org/)
- npm 9+ (comes with Node.js)

# Realtime Queue System

A simple realtime queue system using Go (Gin, WebSocket) for the backend and Vue 3 (Vite) for the frontend.

> **Note:** The queue logic is currently in-memory (data resets on server restart) and fully Dockerized for easy testing and deployment. If you need persistent queue storage, you can extend the backend to use a database (e.g., SQLite, PostgreSQL).

## Features
- REST API and WebSocket for queue management
- Per-client queue assignment (sequential, unique)
- Real-time queue updates and clearing via WebSocket
- Dockerized for easy deployment
- Modern, responsive frontend UI

## Project Structure
```
queue-api/         # Go backend (Gin, WebSocket)
queue-frontend/    # Vue 3 frontend (Vite)
queue-api-demo/    # (Optional) Demo or test code
```

## Quick Start (Docker Compose)
1. Build and run all services:
   ```sh
   docker-compose up --build
   ```
2. Frontend: http://localhost:5173
3. Backend API: http://localhost:8080

## Development
### Backend (Go)
- Location: `queue-api/`
- Run locally:
  ```sh
  cd queue-api
  go run main.go
  ```
- Endpoints:
  - `GET /queue` - Get or assign queue for client
  - `POST /queue/next` - Next queue for client
  - `POST /queue/clear` - Clear all queues (broadcast)
  - `GET /ws` - WebSocket endpoint

### Frontend (Vue 3)
- Location: `queue-frontend/`
- Run locally:
  ```sh
  cd queue-frontend
  npm install
  npm run dev
  ```


## SQLite DB Usage (queue.db)

The backend uses SQLite for persistent queue storage. Here are some common commands:

- **Inspect queue data:**
  ```sh
  sqlite3 queue-api/db/queue.db "SELECT * FROM queue;"
  ```
- **Clear all queue data:**
  ```sh
  sqlite3 queue-api/db/queue.db "DELETE FROM queue;"
  ```
- **Reset DB (delete file and re-init):**
  ```sh
  rm queue-api/db/queue.db
  # Then restart backend to auto-create schema
  ```


If running in Docker, you can run commands directly or enter the container shell:

```
# Run SQL directly from host
docker-compose exec queue-api sqlite3 /app/queue.db "SELECT * FROM queue;"
docker-compose exec queue-api sqlite3 /app/queue.db "DELETE FROM queue;"

# Or enter the container shell for advanced usage
docker-compose exec queue-api sh
# Then inside the container, run:
sqlite3 /app/queue.db
# (You can now use the sqlite3 CLI interactively)
```

---


## How It Works
- Each client gets a unique queue number (A1, A2, ...)
- All queue actions are synced in real-time via WebSocket
- Clearing the queue notifies all clients to reset


## Customization
- Edit UI in `queue-frontend/src/pages/`
- Edit backend logic in `queue-api/main.go`


## License
MIT
