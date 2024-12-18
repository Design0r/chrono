package db

import (
	"context"
	"database/sql"
	"io/fs"
	"log"

	"github.com/pressly/goose/v3"
)

func NewDB(name string) *sql.DB {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
		return nil
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
			log.Printf("Failed to execute PRAGMA statement '%s': %v", p, err)
		}
	}

	RunMigrations(db)

	return db
}

func RunMigrations(db *sql.DB) {
	dir, err := fs.Sub(MigrationFS, "migrations")
	if err != nil {
		log.Fatalf("Failed to find migrations: %v", err)
	}
	prov, err := goose.NewProvider(goose.DialectSQLite3, db, dir)
	if err != nil {
		log.Fatalf("Failed to create goose provider: %v", err)
	}
	if _, err := prov.Up(context.Background()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations ran successfully!")
}
