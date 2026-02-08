package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

// Queue represents a queue entry
 type Queue struct {
	ID         int
	ClientID   string
	QueueNumber string
	CreatedAt  time.Time
	Status     string
}

func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AddQueue(db *sql.DB, clientID, queueNumber string) error {
	_, err := db.Exec(`INSERT INTO queue (client_id, queue_number, created_at, status) VALUES (?, ?, ?, 'active')`, clientID, queueNumber, time.Now().Format("2006-01-02 15:04:05"))
	return err
}

func GetQueueByClientID(db *sql.DB, clientID string) (*Queue, error) {
	row := db.QueryRow(`SELECT id, client_id, queue_number, created_at, status FROM queue WHERE client_id = ? ORDER BY id DESC LIMIT 1`, clientID)
	var q Queue
	var created string
	if err := row.Scan(&q.ID, &q.ClientID, &q.QueueNumber, &created, &q.Status); err != nil {
		return nil, err
	}
	q.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
	return &q, nil
}

func GetLastQueue(db *sql.DB) (*Queue, error) {
	row := db.QueryRow(`SELECT id, client_id, queue_number, created_at, status FROM queue ORDER BY id DESC LIMIT 1`)
	var q Queue
	var created string
	if err := row.Scan(&q.ID, &q.ClientID, &q.QueueNumber, &created, &q.Status); err != nil {
		return nil, err
	}
	q.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
	return &q, nil
}

func ClearQueues(db *sql.DB) error {
	_, err := db.Exec(`UPDATE queue SET status = 'inactive' WHERE status = 'active'`)
	return err
}

// GetQueuesToday returns all queue entries created today
func GetQueuesToday(db *sql.DB) ([]Queue, error) {
	today := time.Now().Format("2006-01-02")
	rows, err := db.Query(`SELECT id, client_id, queue_number, created_at, status FROM queue WHERE DATE(created_at) = ?`, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var queues []Queue
	for rows.Next() {
		var q Queue
		var created string
		if err := rows.Scan(&q.ID, &q.ClientID, &q.QueueNumber, &created, &q.Status); err != nil {
			return nil, err
		}
		q.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", created)
		queues = append(queues, q)
	}
	return queues, nil
}
