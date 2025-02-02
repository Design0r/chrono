package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"

	"chrono/db/repo"
	"chrono/service"
)

// Helper: create sqlmock, wrap in repo.Queries
func initMockDB(t *testing.T) (sqlmock.Sqlmock, *repo.Queries) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	queries := repo.New(db)
	return mock, queries
}

// --------------------------------------------------------
// Test for CreateUser
// --------------------------------------------------------
func TestCreateUser(t *testing.T) {
	mock, queries := initMockDB(t)

	// We expect an INSERT ... RETURNING (sqlc style).
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("testuser", "#FFFFFF", int64(20), "test@example.com", "hashedpw", false).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).
				AddRow(
					int64(1),
					"testuser",
					"test@example.com",
					"hashedpw",
					int64(20),
					false,
					time.Now(),
					time.Now(),
					"#FFFFFF",
				),
		)

	data := repo.CreateUserParams{
		Username:     "testuser",
		Email:        "test@example.com",
		Color:        "#FFFFFF",
		Password:     "hashedpw",
		VacationDays: 20,
		IsSuperuser:  false,
	}

	user, err := service.CreateUser(queries, data)
	if err != nil {
		t.Errorf("CreateUser() returned error: %v", err)
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username to be 'testuser', got %s", user.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for UpdateUser
// --------------------------------------------------------
func TestUpdateUser(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectQuery("UPDATE users").
		WithArgs("#000000", "newName", "new-email@example.com", int64(1234)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).
				AddRow(
					int64(1234),
					"newName",
					"new-email@example.com",
					"hashedpw",
					int64(20),
					false,
					time.Now(),
					time.Now(),
					"#000000",
				),
		)

	params := repo.UpdateUserParams{
		Color:    "#000000",
		Email:    "new-email@example.com",
		Username: "newName",
		ID:       1234,
	}
	updated, err := service.UpdateUser(queries, params)
	if err != nil {
		t.Errorf("UpdateUser() returned error: %v", err)
	}
	if updated.Username != "newName" {
		t.Errorf("Expected updated username to be 'newName', got %s", updated.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for GetUserById
// --------------------------------------------------------
func TestGetUserById(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectQuery(repo.GetUserByID).
		WithArgs(int64(42)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).
				AddRow(
					int64(42),
					"theAnswer",
					"answer@example.com",
					"hashedpw",
					int64(10),
					false,
					time.Now(),
					time.Now(),
					"#FFFFFF",
				),
		)

	u, err := service.GetUserById(queries, 42)
	if err != nil {
		t.Errorf("GetUserById() error = %v", err)
	}
	if u.Username != "theAnswer" {
		t.Errorf("Expected username 'theAnswer', got %s", u.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for GetUserByName
// --------------------------------------------------------
func TestGetUserByName(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectQuery(repo.GetUserByName).
		WithArgs("Alice").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).AddRow(
				int64(11),
				"Alice",
				"alice@example.com",
				"hashedpw",
				int64(20),
				true,
				time.Now(),
				time.Now(),
				"#ABCDEF",
			),
		)

	u, err := service.GetUserByName(queries, "Alice")
	if err != nil {
		t.Errorf("GetUserByName() error = %v", err)
	}
	if u.Username != "Alice" {
		t.Errorf("Expected 'Alice', got %s", u.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for GetUserByEmail
// --------------------------------------------------------
func TestGetUserByEmail(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectQuery(repo.GetUserByEmail).
		WithArgs("bob@example.com").
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).AddRow(
				int64(99),
				"bob",
				"bob@example.com",
				"bobpw",
				int64(15),
				false,
				time.Now(),
				time.Now(),
				"#123456",
			),
		)

	u, err := service.GetUserByEmail(queries, "bob@example.com")
	if err != nil {
		t.Errorf("GetUserByEmail() error = %v", err)
	}
	if u.Username != "bob" {
		t.Errorf("Expected username 'bob', got %s", u.Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for DeleteUser
// --------------------------------------------------------
func TestDeleteUser(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectExec("DELETE FROM users WHERE id").
		WithArgs(int64(77)).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := service.DeleteUser(queries, 77)
	if err != nil {
		t.Errorf("DeleteUser() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for GetCurrentUser
// --------------------------------------------------------
func TestGetCurrentUser(t *testing.T) {
	// This function uses echo.Context to read a cookie, then calls GetUserFromSession
	// internally. We can mock the context and also mock the DB call to GetUserFromSession.

	// 1. Setup echo and request with "session" cookie
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	cookie := &http.Cookie{Name: "session", Value: "some-session-id"}
	req.AddCookie(cookie)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// 2. Setup DB mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating sqlmock: %v", err)
	}
	defer db.Close()

	queries := repo.New(db)

	// Because GetCurrentUser -> GetUserFromSession -> r.GetUserFromSession,
	// we mock that "SELECT ... FROM sessions JOIN users" query:
	mock.ExpectQuery(repo.GetUserFromSession).
		WithArgs("some-session-id"). // second arg is time.Now() in a param
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", // user id
				"username",
				"email",
				"password",
				"vacation_days",
				"is_superuser",
				"created_at",
				"edited_at",
				"color",
			}).AddRow(
				int64(123),
				"sessionUser",
				"sess@example.com",
				"pw",
				int64(10),
				false,
				time.Now(),
				time.Now(),
				"#FFFFFF",
			),
		)

	// 3. Call GetCurrentUser
	user, err := service.GetCurrentUser(queries, ctx)
	if err != nil {
		t.Errorf("GetCurrentUser() error = %v", err)
	}
	if user.Username != "sessionUser" {
		t.Errorf("Expected 'sessionUser', got %s", user.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for GetAllVacUsers
// --------------------------------------------------------
func TestGetAllVacUsers(t *testing.T) {
	mock, queries := initMockDB(t)

	// 1. Mock call to GetUsersWithVacationCount
	mock.ExpectQuery(`(?s)^-- name: GetUsersWithVacationCount :many SELECT u\.id, u\.username, u\.email, u\.password, u\.vacation_days, u\.is_superuser, u\.created_at, u\.edited_at, u\.color, COALESCE\(SUM\(vt.value\), 0\.0\) AS vac_remaining, COALESCE\(SUM\(0\.5\), 0\.0\) AS vac_used FROM users AS u LEFT JOIN vacation_tokens vt ON u\.id = vt\.user_id AND vt\.start_date <= \? AND vt\.end_date >= \? GROUP BY u\.id ORDER BY u\.id`).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
				"vac_remaining", "vac_used", // additional column from your custom view
			}).AddRow(
				int64(10),
				"Alice",
				"alice@example.com",
				"pw",
				int64(20),
				false,
				time.Now(),
				time.Now(),
				"#FF0000",
				10, // your custom query might or might not populate it
				10,
			).AddRow(
				int64(11),
				"Bob",
				"bob@example.com",
				"pw",
				int64(25),
				false,
				time.Now(),
				time.Now(),
				"#00FF00",
				10,
				10,
			),
		)

	// 2. Because of the loop in GetAllVacUsers, each user calls
	//    GetVacationCountForUser for the current year.
	mock.ExpectQuery(repo.GetVacationCountForUser).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(float64(5))) // for Alice
	mock.ExpectQuery(repo.GetVacationCountForUser).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(float64(3))) // for Bob

	users, err := service.GetAllVacUsers(queries)
	if err != nil {
		t.Fatalf("GetAllVacUsers() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	if users[0].VacUsed == nil || users[0].VacUsed != 5 {
		t.Errorf("Expected Alice to have VacUsed=5, got %v", users[0].VacUsed)
	}
	if users[1].VacUsed == nil || users[1].VacUsed != 3 {
		t.Errorf("Expected Bob to have VacUsed=3, got %v", users[1].VacUsed)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for ToggleAdmin
// --------------------------------------------------------
func TestToggleAdmin(t *testing.T) {
	mock, queries := initMockDB(t)

	// The function calls r.ToggleAdmin, then CreateUserNotification
	// 1) toggleAdmin => returns updated user
	mock.ExpectQuery("UPDATE users").
		WithArgs(int64(99)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).AddRow(
				int64(99),
				"targetUser",
				"target@example.com",
				"pw",
				int64(20),
				true, // toggled to admin
				time.Now(),
				time.Now(),
				"#FFFFF0",
			),
		)

	// 2) create notification for userId=99
	mock.ExpectQuery("INSERT INTO notifications").
		WithArgs(sqlmock.AnyArg()). // message
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "message", "created_at", "viewed_at"}).
				AddRow(int64(1), "Some Message", time.Now(), time.Now()),
		)
	// 3) create notification_user association
	mock.ExpectExec("INSERT INTO notification_user").
		WithArgs(int64(1), int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	editor1 := repo.User{ID: 1, Username: "EditorUser", IsSuperuser: true}
	user, err := service.ToggleAdmin(queries, editor1, 99)
	if err != nil {
		t.Errorf("ToggleAdmin() error = %v", err)
	}
	if !user.IsSuperuser {
		t.Errorf("Expected user to be superuser, got false")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for GetAllUsers
// --------------------------------------------------------

func TestGetAllUsers(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectQuery(repo.GetAllUsers).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).
				AddRow(int64(2), "B", "b@example.com", "pw", int64(15), false, time.Now(), time.Now(), "#222222").
				AddRow(int64(3), "C", "c@example.com", "pw", int64(15), false, time.Now(), time.Now(), "#333333"),
		)

	users, err := service.GetAllUsers(queries)
	if err != nil {
		t.Errorf("GetAllUsers() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
	if users[0].Username != "B" || users[1].Username != "C" {
		t.Errorf("Usernames mismatch: got %v, %v", users[0].Username, users[1].Username)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for SetUserVacation
// --------------------------------------------------------
func TestSetUserVacation(t *testing.T) {
	mock, queries := initMockDB(t)

	// 1) GetUserById call
	mock.ExpectQuery(repo.GetUserByID).
		WithArgs(int64(55)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).
				AddRow(
					int64(55),
					"VacUser",
					"vac@example.com",
					"pw",
					int64(15),
					false,
					time.Now(),
					time.Now(),
					"#FFF000",
				),
		)

	// 2) r.SetUserVacation update
	mock.ExpectQuery("UPDATE users").
		WithArgs(int64(30), int64(55)).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"id", "username", "email", "password", "vacation_days",
				"is_superuser", "created_at", "edited_at", "color",
			}).
				AddRow(
					int64(55),
					"VacUser",
					"vac@example.com",
					"pw",
					int64(30),
					false,
					time.Now(),
					time.Now(),
					"#FFF000",
				),
		)

	// add token refresh
	mock.ExpectQuery("INSERT INTO token_refresh").
		WillReturnRows(sqlmock.NewRows([]string{"id", "year", "user_id", "created_at"}).
			AddRow(int64(1), int64(2025), int64(55), time.Now()),
		)

	// create token
	mock.ExpectQuery("INSERT INTO vacation_tokens").
		WillReturnRows(sqlmock.NewRows([]string{"id", "start_sate", "end_date", "value", "user_id"}).
			AddRow(int64(2), time.Now(), time.Now(), 15.0, int64(55)),
		)

	err := service.SetUserVacation(queries, 55, 30, 2025)
	if err != nil {
		t.Errorf("SetUserVacation() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// --------------------------------------------------------
// Test for SetUserColor
// --------------------------------------------------------
func TestSetUserColor(t *testing.T) {
	mock, queries := initMockDB(t)

	mock.ExpectExec("UPDATE users SET color").
		WithArgs("#ABC123", int64(999)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := service.SetUserColor(queries, 999, "#ABC123")
	if err != nil {
		t.Errorf("SetUserColor() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

// Optionally, if you test failing scenarios or errors, you can set up
// mock.Expect*().WillReturnError(...) to ensure you handle DB errors gracefully.
