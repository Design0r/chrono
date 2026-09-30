package db

import (
	"context"
	"database/sql"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

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

func TestRunningTimestampSurvivesNewerFinishedOne(t *testing.T) {
	r, conn := newTestTimestampsRepo(t)
	ctx := context.Background()
	userId := seedUser(t, conn, "a@example.com")

	running, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}

	// an admin adds a finished entry for the user while the timer is running
	_, err = conn.Exec(
		`INSERT INTO timestamps (user_id, start_time, end_time)
		VALUES (?, '2026-01-05 08:00:00', '2026-01-05 16:00:00')`,
		userId,
	)
	if err != nil {
		t.Fatalf("failed to seed finished timestamp: %v", err)
	}

	latest, err := r.GetLatest(ctx, userId)
	if err != nil {
		t.Fatalf("get latest failed: %v", err)
	}
	if latest.ID != running.ID {
		t.Fatalf("expected running timestamp %d as latest, got %d", running.ID, latest.ID)
	}

	again, err := r.Start(ctx, userId)
	if err != nil {
		t.Fatalf("second start failed: %v", err)
	}
	if again.ID != running.ID {
		t.Fatalf("expected running timestamp %d to be reused, got %d", running.ID, again.ID)
	}

	if _, err := r.Stop(ctx, running.ID); err != nil {
		t.Fatalf("stop failed: %v", err)
	}

	latest, err = r.GetLatest(ctx, userId)
	if err != nil {
		t.Fatalf("get latest after stop failed: %v", err)
	}
	if latest.EndTime == nil {
		t.Fatalf("expected no running timestamp after stop, got %+v", latest)
	}
	if count := countTimestamps(t, conn, userId); count != 2 {
		t.Fatalf("expected 2 timestamps, got %d", count)
	}
}

func TestSecondRunningTimestampIsRejected(t *testing.T) {
	r, conn := newTestTimestampsRepo(t)
	ctx := context.Background()
	userId := seedUser(t, conn, "a@example.com")

	if _, err := r.Start(ctx, userId); err != nil {
		t.Fatalf("start failed: %v", err)
	}

	_, err := conn.Exec(`INSERT INTO timestamps (user_id) VALUES (?)`, userId)
	if err == nil {
		t.Fatalf("expected a second running timestamp to be rejected")
	}
}

func TestMigrationRemovesSupersededUnfinishedTimestamps(t *testing.T) {
	const beforeMigration = 20260804122000

	conn, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	dir, err := fs.Sub(chronodb.MigrationFS, "migrations")
	if err != nil {
		t.Fatalf("failed to find migrations: %v", err)
	}
	prov, err := goose.NewProvider(goose.DialectSQLite3, conn, dir)
	if err != nil {
		t.Fatalf("failed to create goose provider: %v", err)
	}

	ctx := context.Background()
	if _, err := prov.UpTo(ctx, beforeMigration); err != nil {
		t.Fatalf("failed to migrate to previous version: %v", err)
	}

	userA := seedUser(t, conn, "a@example.com")
	userB := seedUser(t, conn, "b@example.com")

	// ids: 1 superseded, 2 finished, 3 superseded, 4 running (user a), 5 running (user b)
	seeds := []struct {
		userId   int64
		finished bool
	}{
		{userA, false},
		{userA, true},
		{userA, false},
		{userA, false},
		{userB, false},
	}
	for _, s := range seeds {
		var endTime any
		if s.finished {
			endTime = "2026-01-05 16:00:00"
		}
		_, err := conn.Exec(
			`INSERT INTO timestamps (user_id, start_time, end_time)
			VALUES (?, '2026-01-05 08:00:00', ?)`,
			s.userId, endTime,
		)
		if err != nil {
			t.Fatalf("failed to seed timestamp: %v", err)
		}
	}

	if _, err := prov.Up(ctx); err != nil {
		t.Fatalf("failed to run migration: %v", err)
	}

	rows, err := conn.Query(`SELECT id FROM timestamps ORDER BY id`)
	if err != nil {
		t.Fatalf("failed to read timestamps: %v", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("failed to scan id: %v", err)
		}
		ids = append(ids, id)
	}

	if !slices.Equal(ids, []int64{2, 4, 5}) {
		t.Fatalf("expected timestamps [2 4 5] to remain, got %v", ids)
	}
}
