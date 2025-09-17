package main

import (
	"fmt"
	"log"
	"net/http"

	"antiHate/config"
	"antiHate/internal/bot"
	"antiHate/internal/db"
	"antiHate/internal/httpserver"
	"antiHate/schema/migrations"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	if err := migrations.RunMigrations(database.DB); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("✅ Migrations applied successfully")

	go bot.Start(cfg, database)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("HTTP server running on %s", addr)
	if err := http.ListenAndServe(addr, httpserver.NewRouter(database)); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
