package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type Storage struct {
	mu            sync.RWMutex
	token         string
	createdAt     time.Time
	tokenLifetime time.Duration
}

var storage = &Storage{
	tokenLifetime: 120 * time.Second,
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

func generateHMACToken(data, secretKey, time string) string {
	h := hmac.New(sha256.New, []byte(secretKey+time))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyHMACToken(data, token, secretKey, time string) bool {
	expectedToken := generateHMACToken(data, secretKey, time)
	return hmac.Equal([]byte(token), []byte(expectedToken))
}

// Storage management
func (s *Storage) getOrCreateToken(data, secretKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the current token is still valid
	println(s.token, s.tokenLifetime-time.Since(s.createdAt))
	if s.token != "" && time.Since(s.createdAt) <= s.tokenLifetime {
		return s.token, nil
	}

	// Generate a new token
	s.createdAt = time.Now()
	s.token = generateHMACToken(data, secretKey, s.createdAt.Format(time.RFC3339))

	return s.token, nil
}

func (s *Storage) verifyToken(data, token, secretKey string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if the token matches and is still valid
	if s.token == "" || time.Since(s.createdAt) > s.tokenLifetime {
		return ErrTokenExpired
	}

	if !verifyHMACToken(data, token, secretKey, s.createdAt.Format(time.RFC3339)) {
		return ErrInvalidToken
	}

	return nil
}

// Public API
func GenerateToken(data, secretKey string) (string, error) {
	return storage.getOrCreateToken(data, secretKey)
}

func VerifyToken(data, token, secretKey string) error {
	return storage.verifyToken(data, token, secretKey)
}
