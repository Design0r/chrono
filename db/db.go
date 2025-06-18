package db

import (
	"context"
	"database/sql"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
)

func NewDB(name string) *sql.DB {
	err := os.MkdirAll("db", 0755)
	if err != nil {
		slog.Error("Failed to create directory", "error", err)
	}

	dbPath := filepath.Join("db", name)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("Failed to open database", "error", err)
	}

	pragmas := []string{
		"PRAGMA foreign_keys=ON;",
		"PRAGMA journal_mode=WAL;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA temp_store=MEMORY;",
		"PRAGMA mmap_size=134217728;",
		"PRAGMA journal_size_limit=67108864;",
		"PRAGMA cache_size=2000;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			slog.Error("Failed to execute PRAGMA statement", "pragma", p, "error", err)
		}
	}

	RunMigrations(db)

	return db
}

func CloseDB(db *sql.DB) {
	slog.Info("Closing Database")
	db.Close()
}

func RunMigrations(db *sql.DB) {
	dir, err := fs.Sub(MigrationFS, "migrations")
	if err != nil {
		slog.Error("Failed to find migrations", "error", err)
		os.Exit(1)
	}
	prov, err := goose.NewProvider(goose.DialectSQLite3, db, dir)
	if err != nil {
		slog.Error("Failed to create goose provider", "error", err)
		os.Exit(1)
	}
	if _, err := prov.Up(context.Background()); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	slog.Info("Migrations ran successfully!")
}
