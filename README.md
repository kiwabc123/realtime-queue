# Requirements

## Backend
- Go 1.20+ (https://golang.org/dl/)
- Docker & Docker Compose (for containerized run)

## Frontend
- Node.js 18+ (https://nodejs.org/)
- npm 9+ (comes with Node.js)

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

## How It Works
- Each client gets a unique queue number (A1, A2, ...)
- All queue actions are synced in real-time via WebSocket
- Clearing the queue notifies all clients to reset

## Customization
- Edit UI in `queue-frontend/src/pages/`
- Edit backend logic in `queue-api/main.go`

## License
MIT
