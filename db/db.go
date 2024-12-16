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
