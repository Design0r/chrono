package db

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	chronodb "chrono/db"
	"chrono/db/repo"
	"chrono/internal/domain"
)

func newTestTimestampsRepo(t *testing.T) (domain.TimestampsRepository, *sql.DB) {
	t.Helper()

	conn, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	chronodb.RunMigrations(conn)

	log := slog.New(slog.NewTextHandler(io.Discard, nil))

	return NewSQLTimestampsRepo(repo.New(conn), log), conn
}

func seedUser(t *testing.T, conn *sql.DB, email string) int64 {
	t.Helper()

	res, err := conn.Exec(
		`INSERT INTO users (username, email, password, vacation_days, is_superuser)
		VALUES (?, ?, '', 30, 0)`,
		email, email,
	)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		t.Fatalf("failed to read user id: %v", err)
	}

	return id
}

func countTimestamps(t *testing.T, conn *sql.DB, userId int64) int {
	t.Helper()

	var count int
	err := conn.QueryRow(`SELECT COUNT(*) FROM timestamps WHERE user_id = ?`, userId).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count timestamps: %v", err)
	}

	return count
}

func TestStartReusesRunningTimestamp(t *testing.T) {
	r, conn := newTestTimestampsRepo(t)
	ctx := context.Background()
	userId := seedUser(t, conn, "a@example.com")

	first, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("first start failed: %v", err)
	}

	// second device with outdated state starts a timer again
	second, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("second start failed: %v", err)
	}

	if second.ID != first.ID {
		t.Fatalf("expected running timestamp %d to be reused, got %d", first.ID, second.ID)
	}
	if second.EndTime != nil {
		t.Fatalf("expected reused timestamp to be still running")
	}
	if !second.StartTime.Equal(first.StartTime) {
		t.Fatalf("expected start time %v to be kept, got %v", first.StartTime, second.StartTime)
	}
	if count := countTimestamps(t, conn, userId); count != 1 {
		t.Fatalf("expected 1 timestamp, got %d", count)
	}
}

func TestStartCreatesNewTimestampAfterStop(t *testing.T) {
	r, conn := newTestTimestampsRepo(t)
	ctx := context.Background()
	userId := seedUser(t, conn, "a@example.com")

	first, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("first start failed: %v", err)
	}
	if _, err := r.Stop(ctx, first.ID); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	second, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("second start failed: %v", err)
	}

	if second.ID == first.ID {
		t.Fatalf("expected a new timestamp after the first one was stopped")
	}
	if second.EndTime != nil {
		t.Fatalf("expected new timestamp to be running")
	}
	if count := countTimestamps(t, conn, userId); count != 2 {
		t.Fatalf("expected 2 timestamps, got %d", count)
	}
}

func TestStartIgnoresRunningTimestampOfOtherUser(t *testing.T) {
	r, conn := newTestTimestampsRepo(t)
	ctx := context.Background()
	userA := seedUser(t, conn, "a@example.com")
	userB := seedUser(t, conn, "b@example.com")

	a, err := r.Start(ctx, userA)
	if err != nil {
		t.Fatalf("start for user a failed: %v", err)
	}

	b, err := r.Start(ctx, userB)
	if err != nil {
		t.Fatalf("start for user b failed: %v", err)
	}

	if b.ID == a.ID || b.UserID != userB {
		t.Fatalf("expected user b to get an own timestamp, got %+v", b)
	}
}

func TestStartIgnoresUnfinishedTimestampBehindNewerOne(t *testing.T) {
	r, conn := newTestTimestampsRepo(t)
	ctx := context.Background()
	userId := seedUser(t, conn, "a@example.com")

	// leftover from before the check existed: never ended, but superseded
	_, err := conn.Exec(`INSERT INTO timestamps (user_id) VALUES (?)`, userId)
	if err != nil {
		t.Fatalf("failed to seed unfinished timestamp: %v", err)
	}
	_, err = conn.Exec(
		`INSERT INTO timestamps (user_id, end_time) VALUES (?, CURRENT_TIMESTAMP)`,
		userId,
	)
	if err != nil {
		t.Fatalf("failed to seed finished timestamp: %v", err)
	}

	ts, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	latest, err := r.GetLatest(ctx, userId)
	if err != nil {
		t.Fatalf("get latest failed: %v", err)
	}

	if ts.ID != latest.ID || ts.EndTime != nil {
		t.Fatalf("expected a new running timestamp as latest, got %+v (latest %+v)", ts, latest)
	}
	if count := countTimestamps(t, conn, userId); count != 3 {
		t.Fatalf("expected 3 timestamps, got %d", count)
	}
}
