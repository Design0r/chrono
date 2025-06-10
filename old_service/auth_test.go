package service_test

import (
	"testing"
	"time"

	"chrono/db/repo"
	"chrono/service"
)

// TestHashAndCheckPassword ensures a hashed password matches the original.
func TestHashAndCheckPassword(t *testing.T) {
	password := "mySuperSecret123!"
	hashed, err := service.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error hashing password, got %v", err)
	}

	if hashed == password {
		t.Errorf("Hash should not be the same as plain text.")
	}

	if !service.CheckPassword(hashed, password) {
		t.Errorf("Expected CheckPassword to return true for correct password.")
	}

	if service.CheckPassword(hashed, "wrongPassword") {
		t.Errorf("Expected CheckPassword to return false for incorrect password.")
	}
}

// TestCreateSessionCookie verifies the structure of the cookie returned.
func TestCreateSessionCookie(t *testing.T) {
	session := repo.Session{
		ID:         "randomSessionID",
		ValidUntil: time.Now().Add(24 * time.Hour),
	}

	cookie := service.CreateSessionCookie(session)
	if cookie.Name != "session" {
		t.Errorf("Expected cookie name 'session', got %s", cookie.Name)
	}
	if cookie.Value != session.ID {
		t.Errorf("Expected cookie value %s, got %s", session.ID, cookie.Value)
	}
	if cookie.Expires != session.ValidUntil {
		t.Errorf("Expected cookie expires %v, got %v", session.ValidUntil, cookie.Expires)
	}
	if cookie.HttpOnly != true {
		t.Errorf("Expected cookie HttpOnly = true")
	}
}

// TestDeleteSessionCookie verifies that the "deletion" cookie is formed properly.
func TestDeleteSessionCookie(t *testing.T) {
	cookie := service.DeleteSessionCookie()
	if cookie.Value != "" {
		t.Errorf("Expected empty cookie value, got %s", cookie.Value)
	}
	if cookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge = -1 for deletion, got %d", cookie.MaxAge)
	}
	if cookie.Expires.After(time.Now()) {
		t.Errorf("Expected Expires to be in the past, got %v", cookie.Expires)
	}
	if cookie.HttpOnly != true {
		t.Errorf("Expected HttpOnly = true")
	}
}

// TestSecureRandom checks that returned string has the requested length (in base64).
func TestSecureRandom(t *testing.T) {
	length := 32
	randomString := service.SecureRandom(length)
	if len(randomString) == 0 {
		t.Errorf("Expected non-empty random string.")
	}
	// Because it's base64-encoded, the length will differ from the raw byte length.
	// This check is just a sanity check that we got something roughly correct.
	if len(randomString) < length {
		t.Errorf("Expected at least length %d, got %d", length, len(randomString))
	}
}
