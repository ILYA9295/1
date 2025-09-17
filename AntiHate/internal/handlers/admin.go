package handlers

import (
	"antiHate/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterAdminRoutes(r chi.Router, repo *repository.Repository) {
	r.Get("/admin/messages", func(w http.ResponseWriter, r *http.Request) {
		rows, err := repo.GetAllMessages()
		if err != nil {
			http.Error(w, "DB error", 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rows)
	})
}

//не нужен он
