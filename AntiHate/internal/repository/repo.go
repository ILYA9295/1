package repository

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

type Message struct {
	ID        int       `db:"id" json:"id"`
	ChatID    int64     `db:"chat_id" json:"chat_id"`
	UserID    int64     `db:"user_id" json:"user_id"`
	Message   string    `db:"message" json:"message"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

func New(url string) (*Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

func (r *Repository) Close() {
	r.db.Close()
}

func (r *Repository) SaveMessage(chatID int64, userID int64, message string) error {
	_, err := r.db.Exec(
		"INSERT INTO messages (chat_id, user_id, message, created_at) VALUES ($1,$2,$3,NOW())",
		chatID, userID, message,
	)
	return err
}

// Новый метод для получения всех сообщений
func (r *Repository) GetAllMessages() ([]Message, error) {
	rows, err := r.db.Query("SELECT id, chat_id, user_id, message, created_at FROM messages ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ChatID, &m.UserID, &m.Message, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
