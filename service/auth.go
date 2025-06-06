package service

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"chrono/config"
	"chrono/db/repo"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Returns true if the password matches
}

func CreateSessionCookie(session repo.Session) *http.Cookie {
	return &http.Cookie{
		Path:     "/",
		Name:     "session",
		Value:    session.ID,
		HttpOnly: true,
		Secure:   !config.GetConfig().Debug, // true if debug off
		Expires:  session.ValidUntil,
		SameSite: http.SameSiteStrictMode,
	}
}

func DeleteSessionCookie() *http.Cookie {
	return &http.Cookie{
		Path:     "/",
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Secure:   !config.GetConfig().Debug, // true if debug off
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}
}

func SecureRandom(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}
