# Queue DB Integration

This directory contains the schema and Go code for SQLite-based queue storage.

## Schema
- `schema.sql`: Defines tables for queue and queue_log.

## Go Code
- `queue.go`: Provides functions to initialize DB, add queue, get queue by client, get last queue, and clear queues.

## Setup
1. Install dependencies:
   ```sh
   cd queue-api
   go mod tidy
   ```
2. Initialize DB:
   ```sh
   sqlite3 queue.db < db/schema.sql
   ```
3. Use functions in `db/queue.go` to interact with the queue.

## Extend
- Add more fields or tables as needed for your use case.
- Integrate with main.go for persistent queue logic.
