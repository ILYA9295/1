package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"strings"
)

//go:embed 001_init.sql
var migrationFiles embed.FS

func RunMigrations(db *sql.DB) error {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		data, err := migrationFiles.ReadFile(name)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", name, err)
		}

		sqlQuery := string(data)
		if strings.TrimSpace(sqlQuery) == "" {
			continue
		}

		log.Printf("🚀 Running migration: %s", name)
		if _, err := db.Exec(sqlQuery); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", name, err)
		}
	}
	return nil
}
