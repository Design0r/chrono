package db

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"path/filepath"
	"testing"
	"time"

	chronodb "chrono/db"
	"chrono/db/repo"
	"chrono/internal/domain"
)

func newTestSessionRepo(t *testing.T) (domain.SessionRepository, *sql.DB) {
	t.Helper()

	conn, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	chronodb.RunMigrations(conn)

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewSQLSessionRepo(repo.New(conn), log), conn
}

func TestDeleteExpiredSessions(t *testing.T) {
	r, conn := newTestSessionRepo(t)
	ctx := context.Background()
	userId := seedUser(t, conn, "a@example.com")

	now := time.Now()
	berlin := time.FixedZone("CEST", 2*60*60)
	newYork := time.FixedZone("EDT", -4*60*60)

	sessions := map[string]time.Time{
		"expired-utc":    now.Add(-time.Minute).UTC(),
		"expired-berlin": now.Add(-time.Minute).In(berlin),
		// lexically "later" than now in UTC, but already expired
		"expired-new-york": now.Add(-time.Minute).In(newYork),
		"valid-utc":        now.Add(time.Hour).UTC(),
		// lexically "earlier" than now in UTC, but still valid
		"valid-berlin": now.Add(time.Hour).In(berlin),
	}
	for id, validUntil := range sessions {
		_, err := conn.Exec(
			`INSERT INTO sessions (id, valid_until, user_id) VALUES (?, ?, ?)`,
			id, validUntil, userId,
		)
		if err != nil {
			t.Fatalf("failed to seed session %v: %v", id, err)
		}
	}

	deleted, err := r.DeleteExpired(ctx)
	if err != nil {
		t.Fatalf("delete expired failed: %v", err)
	}
	if deleted != 3 {
		t.Fatalf("expected 3 deleted sessions, got %d", deleted)
	}

	for _, id := range []string{"valid-utc", "valid-berlin"} {
		if _, err := r.GetById(ctx, id); err != nil {
			t.Fatalf("expected session %v to be kept: %v", id, err)
		}
	}
	for _, id := range []string{"expired-utc", "expired-berlin", "expired-new-york"} {
		if _, err := r.GetById(ctx, id); err == nil {
			t.Fatalf("expected session %v to be deleted", id)
		}
	}
}
