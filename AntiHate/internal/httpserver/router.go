package httpserver

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Get("/messages", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, chat_id, user_id, message, deleted, created_at FROM messages ORDER BY created_at DESC")
		if err != nil {
			http.Error(w, "DB error", 500)
			return
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
			rows.Scan(&id, &chatID, &userID, &message, &deleted, &createdAt)
			result = append(result, map[string]interface{}{
				"id":         id,
				"chat_id":    chatID,
				"user_id":    userID,
				"message":    message,
				"deleted":    deleted,
				"created_at": createdAt,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	return r
}
