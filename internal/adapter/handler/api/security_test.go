package api_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	chronodb "chrono/db"
	"chrono/db/repo"
	"chrono/internal/adapter/db"
	"chrono/internal/adapter/handler/api"
	mw "chrono/internal/adapter/middleware"
	"chrono/internal/domain"
	"chrono/internal/service"
	"chrono/internal/service/auth"
)

const testPassword = "correct horse battery staple"

type testApp struct {
	e    *echo.Echo
	conn *sql.DB
	pw   auth.PasswordHasher
}

// newTestApp wires the real repositories, services, handlers and middleware
// the same way the server does, backed by a fresh sqlite database.
func newTestApp(t *testing.T) *testApp {
	t.Helper()

	conn, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	chronodb.RunMigrations(conn)

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	q := repo.New(conn)

	userRepo := db.NewSQLUserRepo(q, log)
	sessionRepo := db.NewSQLSessionRepo(q, log)

	tokenSvc := service.NewTokenService(
		db.NewSQLRefreshTokenRepo(q, log),
		db.NewSQLVacationTokenRepo(q, log),
		log,
	)
	notifSvc := service.NewNotificationService(
		db.NewSQLNotificationRepo(q, log),
		db.NewSQLUserNotificationRepo(q, log),
		log,
	)
	userSvc := service.NewUserService(userRepo, notifSvc, tokenSvc, log)
	requestSvc := service.NewRequestService(
		db.NewSQLRequestRepo(q, log),
		userRepo,
		notifSvc,
		log,
	)
	eventSvc := service.NewEventService(
		db.NewSQLEventUserRepo(q, log),
		requestSvc,
		userSvc,
		tokenSvc,
		log,
	)
	pw := auth.NewBcryptHasher(4)
	authSvc := service.NewAuthService(userRepo, sessionRepo, time.Hour, false, pw, log)
	settingsSvc := service.NewSettingsService(db.NewSQLSettingsRepo(q, log), log)
	aworkSvc := service.NewAworkService(eventSvc, userSvc, log)
	timestampsSvc := service.NewTimestampsService(
		db.NewSQLTimestampsRepo(q, log),
		eventSvc,
		log,
		userSvc,
	)

	e := echo.New()
	apiGrp := e.Group("/api/v1")
	authGrp := apiGrp.Group(
		"",
		mw.SessionMiddleware(authSvc),
		mw.AuthenticationMiddleware(authSvc),
	)
	adminGrp := authGrp.Group("", mw.AdminMiddleware())

	authHandler := api.NewAPIAuthHandler(userSvc, authSvc, settingsSvc, log)
	userHandler := api.NewAPIUserHandler(userSvc, eventSvc, authSvc, tokenSvc, log)
	aworkHandler := api.NewAPIAworkHandler(userSvc, eventSvc, aworkSvc, log)
	timestampsHandler := api.NewAPITimestampsHandler(timestampsSvc, userSvc)

	authHandler.RegisterRoutes(apiGrp)
	userHandler.RegisterRoutes(authGrp)
	aworkHandler.RegisterRoutes(authGrp, adminGrp)
	timestampsHandler.RegisterRoutes(authGrp, adminGrp)

	return &testApp{e: e, conn: conn, pw: pw}
}

func (a *testApp) seedUser(t *testing.T, email string, admin bool) int64 {
	t.Helper()

	hash, err := a.pw.Hash(testPassword)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	role := "user"
	if admin {
		role = "admin"
	}

	res, err := a.conn.Exec(
		`INSERT INTO users (username, email, password, vacation_days, is_superuser, role)
		VALUES (?, ?, ?, 30, ?, ?)`,
		email, email, hash, admin, role,
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

func (a *testApp) do(
	t *testing.T,
	method string,
	path string,
	session string,
	form url.Values,
) *httptest.ResponseRecorder {
	t.Helper()

	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}

	req := httptest.NewRequest(method, path, body)
	if form != nil {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	}
	if session != "" {
		req.AddCookie(&http.Cookie{Name: "session", Value: session})
	}

	rec := httptest.NewRecorder()
	a.e.ServeHTTP(rec, req)

	return rec
}

func (a *testApp) login(t *testing.T, email string, password string) (string, int) {
	t.Helper()

	rec := a.do(t, http.MethodPost, "/api/v1/login", "", url.Values{
		"email":    {email},
		"password": {password},
	})

	for _, c := range rec.Result().Cookies() {
		if c.Name == "session" && c.Value != "" {
			return c.Value, rec.Code
		}
	}

	return "", rec.Code
}

func (a *testApp) mustLogin(t *testing.T, email string) string {
	t.Helper()

	session, code := a.login(t, email, testPassword)
	if code != http.StatusOK || session == "" {
		t.Fatalf("login for %v failed with status %d", email, code)
	}

	return session
}

func (a *testApp) sessionCount(t *testing.T, userId int64) int {
	t.Helper()

	var count int
	err := a.conn.QueryRow(`SELECT COUNT(*) FROM sessions WHERE user_id = ?`, userId).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count sessions: %v", err)
	}

	return count
}

func expectStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("expected status %d, got %d: %s", want, rec.Code, rec.Body.String())
	}
}

func TestUserCannotEditOtherUser(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	adminId := app.seedUser(t, "admin2@example.com", true)
	app.seedUser(t, "user@example.com", false)
	user := app.mustLogin(t, "user@example.com")

	rec := app.do(t, http.MethodPatch, fmt.Sprintf("/api/v1/users/%d", adminId), user, url.Values{
		"password": {"taken over"},
	})
	expectStatus(t, rec, http.StatusForbidden)

	if _, code := app.login(t, "admin2@example.com", "taken over"); code == http.StatusOK {
		t.Fatalf("expected the admin password to be unchanged")
	}
}

func TestStopForeignTimestampIsRejected(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	app.seedUser(t, "a@example.com", false)
	app.seedUser(t, "b@example.com", false)
	a := app.mustLogin(t, "a@example.com")
	b := app.mustLogin(t, "b@example.com")

	expectStatus(t, app.do(t, http.MethodPost, "/api/v1/timestamps", a, nil), http.StatusOK)

	var id int64
	if err := app.conn.QueryRow(`SELECT id FROM timestamps`).Scan(&id); err != nil {
		t.Fatalf("failed to read timestamp: %v", err)
	}
	path := fmt.Sprintf("/api/v1/timestamps/%d", id)

	expectStatus(t, app.do(t, http.MethodPatch, path, b, nil), http.StatusNotFound)

	var running bool
	err := app.conn.QueryRow(`SELECT end_time IS NULL FROM timestamps WHERE id = ?`, id).
		Scan(&running)
	if err != nil {
		t.Fatalf("failed to read timestamp: %v", err)
	}
	if !running {
		t.Fatalf("expected timestamp of user a to still be running")
	}

	expectStatus(t, app.do(t, http.MethodPatch, path, a, nil), http.StatusOK)
}

func TestDisabledUserCannotLogin(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	userId := app.seedUser(t, "user@example.com", false)

	if _, err := app.conn.Exec(`UPDATE users SET enabled = 0 WHERE id = ?`, userId); err != nil {
		t.Fatalf("failed to disable user: %v", err)
	}

	session, code := app.login(t, "user@example.com", testPassword)
	if code != http.StatusForbidden || session != "" {
		t.Fatalf("expected login to be rejected, got status %d", code)
	}
}

func TestDisablingUserEndsSessions(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	userId := app.seedUser(t, "user@example.com", false)
	admin := app.mustLogin(t, "admin@example.com")
	user := app.mustLogin(t, "user@example.com")

	expectStatus(t, app.do(t, http.MethodGet, "/api/v1/timestamps/day", user, nil), http.StatusOK)

	rec := app.do(t, http.MethodPatch, fmt.Sprintf("/api/v1/users/%d", userId), admin, url.Values{
		"enabled": {"false"},
	})
	expectStatus(t, rec, http.StatusOK)

	if count := app.sessionCount(t, userId); count != 0 {
		t.Fatalf("expected all sessions of the disabled user to be deleted, got %d", count)
	}
	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/day", user, nil),
		http.StatusUnauthorized,
	)
}

func TestDisabledUserWithLeftoverSessionIsRejected(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	userId := app.seedUser(t, "user@example.com", false)
	user := app.mustLogin(t, "user@example.com")

	// disabled directly in the database, so the session is not cleaned up
	if _, err := app.conn.Exec(`UPDATE users SET enabled = 0 WHERE id = ?`, userId); err != nil {
		t.Fatalf("failed to disable user: %v", err)
	}

	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/day", user, nil),
		http.StatusUnauthorized,
	)
}

func TestPasswordChangeEndsOtherSessions(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	userId := app.seedUser(t, "user@example.com", false)
	current := app.mustLogin(t, "user@example.com")
	other := app.mustLogin(t, "user@example.com")

	rec := app.do(t, http.MethodPatch, fmt.Sprintf("/api/v1/users/%d", userId), current, url.Values{
		"password": {"a new password"},
	})
	expectStatus(t, rec, http.StatusOK)

	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/day", other, nil),
		http.StatusUnauthorized,
	)
	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/day", current, nil),
		http.StatusOK,
	)
	if count := app.sessionCount(t, userId); count != 1 {
		t.Fatalf("expected only the current session to remain, got %d", count)
	}
}

func TestAdminPasswordResetEndsAllSessions(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	userId := app.seedUser(t, "user@example.com", false)
	admin := app.mustLogin(t, "admin@example.com")
	app.mustLogin(t, "user@example.com")

	rec := app.do(t, http.MethodPatch, fmt.Sprintf("/api/v1/users/%d", userId), admin, url.Values{
		"password": {"a new password"},
	})
	expectStatus(t, rec, http.StatusOK)

	if count := app.sessionCount(t, userId); count != 0 {
		t.Fatalf("expected all sessions of the user to be deleted, got %d", count)
	}
}

func TestAworkIDCanOnlyBeChangedByAdmins(t *testing.T) {
	const aworkId = "0b5b2f1e-3c4d-4e5f-8a9b-0c1d2e3f4a5b"

	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	userId := app.seedUser(t, "user@example.com", false)
	admin := app.mustLogin(t, "admin@example.com")
	user := app.mustLogin(t, "user@example.com")
	path := fmt.Sprintf("/api/v1/users/%d", userId)

	rec := app.do(t, http.MethodPatch, path, user, url.Values{"awork_id": {aworkId}})
	expectStatus(t, rec, http.StatusForbidden)

	rec = app.do(t, http.MethodPatch, path, admin, url.Values{
		"awork_id": {"x' or userId ne guid'00000000-0000-0000-0000-000000000000"},
	})
	expectStatus(t, rec, http.StatusUnprocessableEntity)

	rec = app.do(t, http.MethodPatch, path, admin, url.Values{"awork_id": {aworkId}})
	expectStatus(t, rec, http.StatusOK)

	// the profile form of the user resends the unchanged id
	rec = app.do(t, http.MethodPatch, path, user, url.Values{
		"awork_id": {aworkId},
		"username": {"renamed"},
	})
	expectStatus(t, rec, http.StatusOK)

	// the team table does not send the awork id at all, it must stay linked
	rec = app.do(t, http.MethodPatch, path, admin, url.Values{"role": {"user"}})
	expectStatus(t, rec, http.StatusOK)

	var stored sql.NullString
	err := app.conn.QueryRow(`SELECT awork_id FROM users WHERE id = ?`, userId).Scan(&stored)
	if err != nil {
		t.Fatalf("failed to read awork id: %v", err)
	}
	if stored.String != aworkId {
		t.Fatalf("expected awork id %v to be kept, got %q", aworkId, stored.String)
	}
}

func TestAworkDirectoryIsAdminOnly(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	app.seedUser(t, "user@example.com", false)
	user := app.mustLogin(t, "user@example.com")

	expectStatus(t, app.do(t, http.MethodGet, "/api/v1/awork/users", user, nil), http.StatusForbidden)
}

func TestAworkTimeEntriesRejectInvalidID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := service.NewAworkService(nil, nil, log)

	_, err := svc.GetTimeEntries(
		context.Background(),
		"x' or userId ne guid'00000000-0000-0000-0000-000000000000",
		time.Now(),
		time.Now(),
	)
	if err != service.ErrInvalidAworkID {
		t.Fatalf("expected ErrInvalidAworkID, got %v", err)
	}
}

func TestTimestampSummaryIsAdminOnly(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	app.seedUser(t, "user@example.com", false)
	user := app.mustLogin(t, "user@example.com")

	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/all/summary", user, nil),
		http.StatusForbidden,
	)
	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/all?userId=1", user, nil),
		http.StatusForbidden,
	)
}

func TestAdminTimestampsCanBeFilteredByUser(t *testing.T) {
	app := newTestApp(t)
	app.seedUser(t, "admin@example.com", true)
	a := app.seedUser(t, "a@example.com", false)
	b := app.seedUser(t, "b@example.com", false)
	admin := app.mustLogin(t, "admin@example.com")

	createTimestamp := func(userId int64, start string, end string) {
		form := url.Values{
			"user_id":    {fmt.Sprint(userId)},
			"start_time": {start},
			"end_time":   {end},
		}
		expectStatus(
			t,
			app.do(t, http.MethodPost, "/api/v1/timestamps/new", admin, form),
			http.StatusOK,
		)
	}
	createTimestamp(a, "2026-03-02T08:00:00Z", "2026-03-02T10:00:00Z")
	createTimestamp(a, "2026-03-03T08:00:00Z", "2026-03-03T09:30:00Z")
	createTimestamp(b, "2026-03-02T08:00:00Z", "2026-03-02T09:00:00Z")

	const query = "startDate=2026-01-01&endDate=2027-01-01"

	rec := app.do(t, http.MethodGet, fmt.Sprintf("/api/v1/timestamps/all?%s&userId=%d", query, a), admin, nil)
	expectStatus(t, rec, http.StatusOK)
	var filtered struct {
		Data []domain.Timestamp `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &filtered); err != nil {
		t.Fatalf("failed to decode timestamps: %v", err)
	}
	if len(filtered.Data) != 2 {
		t.Fatalf("expected 2 timestamps for user a, got %d", len(filtered.Data))
	}
	for _, ts := range filtered.Data {
		if ts.UserID != a {
			t.Fatalf("expected only timestamps of user %d, got one of user %d", a, ts.UserID)
		}
	}

	rec = app.do(t, http.MethodGet, "/api/v1/timestamps/all/summary?"+query, admin, nil)
	expectStatus(t, rec, http.StatusOK)
	var summary struct {
		Data map[string]float64 `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &summary); err != nil {
		t.Fatalf("failed to decode summary: %v", err)
	}
	want := map[string]float64{fmt.Sprint(a): 3.5 * 3600, fmt.Sprint(b): 3600}
	for id, seconds := range want {
		if summary.Data[id] != seconds {
			t.Fatalf("expected %v worked seconds for user %s, got %v", seconds, id, summary.Data[id])
		}
	}
	if len(summary.Data) != len(want) {
		t.Fatalf("expected summary for %d users, got %v", len(want), summary.Data)
	}

	expectStatus(
		t,
		app.do(t, http.MethodGet, "/api/v1/timestamps/all?userId=abc", admin, nil),
		http.StatusUnprocessableEntity,
	)
}
