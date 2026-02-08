-- SQLite schema for queue system

CREATE TABLE queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id TEXT NOT NULL,
    queue_number TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'active'
);

CREATE TABLE queue_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    queue_id INTEGER,
    action TEXT,
    action_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    note TEXT,
    FOREIGN KEY(queue_id) REFERENCES queue(id)
);