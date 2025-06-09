package auth

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Hash(pw string) (string, error)
	Compare(hashedPw, pw string) bool
	SecureRandom(length int) string
	SecureRandom64() string
}

type bcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) bcryptHasher {
	return bcryptHasher{cost: cost}
}

func (svc *bcryptHasher) SecureRandom(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes)
}

func (svc *bcryptHasher) SecureRandom64() string {
	return svc.SecureRandom(64)
}

func (svc *bcryptHasher) Hash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), svc.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (svc *bcryptHasher) Compare(hashedPw, pw string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw))
	return err == nil // Returns true if the password matches
}
