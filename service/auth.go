package service

import (
	"encoding/base64"
	"net/http"
	"time"

	"crypto/rand"

	"golang.org/x/crypto/bcrypt"

	"calendar/db/repo"
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
	cookie := http.Cookie{}
	cookie.Path = "/"
	cookie.Name = "session"
	cookie.Value = session.ID
	cookie.HttpOnly = true
	cookie.Secure = false // Ensure this is set to true in production
	cookie.Expires = session.ValidUntil

	return &cookie
}

func DeleteSessionCookie() *http.Cookie {
	cookie := http.Cookie{}
	cookie.Path = "/"
	cookie.Name = "session"
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.Secure = false // Ensure this is set to true in production
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1

	return &cookie
}

func SecureRandom(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}
