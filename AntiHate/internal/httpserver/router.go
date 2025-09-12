package httpserver

import (
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
		w.Write([]byte("[]")) // пока заглушка
	})

	return r
}
