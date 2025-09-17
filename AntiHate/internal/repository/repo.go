package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Close() {
	r.db.Close()
}

// Чистое сообщение
func (r *Repository) SaveMessage(chatID int64, userID int64, message string) error {
	_, err := r.db.Exec(
		"INSERT INTO messages (chat_id, user_id, message, deleted, created_at) VALUES ($1,$2,$3,$4,NOW())",
		chatID, userID, message, false,
	)
	return err
}

// Удалённое сообщение
func (r *Repository) SaveDeletedMessage(chatID int64, userID int64, message string) error {
	_, err := r.db.Exec(
		"INSERT INTO messages (chat_id, user_id, message, deleted, created_at) VALUES ($1,$2,$3,$4,NOW())",
		chatID, userID, message, true,
	)
	return err
}

// Получение всех сообщений
func (r *Repository) GetAllMessages() ([]map[string]interface{}, error) {
	rows, err := r.db.Query("SELECT id, chat_id, user_id, message, deleted, created_at FROM messages ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []map[string]interface{}
	for rows.Next() {
		var (
			id        int
			chatID    int64
			userID    int64
			message   string
			deleted   bool
			createdAt string
		)
		if err := rows.Scan(&id, &chatID, &userID, &message, &deleted, &createdAt); err != nil {
			return nil, err
		}

		row := map[string]interface{}{
			"id":         id,
			"chat_id":    chatID,
			"user_id":    userID,
			"message":    message,
			"deleted":    deleted,
			"created_at": createdAt,
		}
		result = append(result, row)
	}
	return result, nil
}
