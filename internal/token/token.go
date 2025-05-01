package token

import (
	"anon-chat/internal/pow"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type RotatingToken struct {
	mu        sync.RWMutex
	value     string
	expiresAt time.Time
}

var ephemeralToken = &RotatingToken{}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

func (s *RotatingToken) getOrCreateToken(userKey, secretKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the token is expired or not set
	if s.value == "" || s.expiresAt.Before(time.Now()) {
		// Generate a new RotatingToken
		token, err := pow.RandomKey()
		if err != nil {
			return "", err
		}

		s.expiresAt = time.Now().Add(2 * time.Minute) // Set token lifetime to 2 minutes
		s.value = token
	}

	return generateHMACToken(
		s.value+"|"+userKey+"|"+s.expiresAt.Format(time.RFC3339),
		secretKey,
	), nil
}

func (s *RotatingToken) verifyToken(userKey, token, secretKey string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if the token matches and is still valid
	if s.value == "" || s.expiresAt.Before(time.Now()) {
		return ErrTokenExpired
	}

	isVerified := verifyHMACToken(
		s.value+"|"+userKey+"|"+s.expiresAt.Format(time.RFC3339),
		token,
		secretKey,
	)

	if !isVerified {
		return ErrInvalidToken
	}

	return nil
}

func generateHMACToken(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyHMACToken(data, token, secretKey string) bool {
	expectedToken := generateHMACToken(data, secretKey)
	return hmac.Equal([]byte(token), []byte(expectedToken))
}

// Public API
func GenerateToken(userKey, secretKey string) (string, error) {
	return ephemeralToken.getOrCreateToken(userKey, secretKey)
}

func VerifyToken(userKey, token, secretKey string) error {
	return ephemeralToken.verifyToken(userKey, token, secretKey)
}
