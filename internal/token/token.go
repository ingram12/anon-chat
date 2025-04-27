package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

const (
	maxTokens = 100
)

type tokenInfo struct {
	createdAt time.Time
}

type TokenStorage struct {
	mu            sync.RWMutex
	tokens        map[string]tokenInfo
	tokenLifetime time.Duration
}

var storage = &TokenStorage{
	tokens:        make(map[string]tokenInfo, maxTokens),
	tokenLifetime: 30 * time.Second,
}

var (
	ErrInvalidToken  = errors.New("invalid token")
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
	ErrStorageFull   = errors.New("token storage is full, please try again later")
)

// Token generation and verification
func generateHMACToken(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyHMACToken(data, token, secretKey string) bool {
	expectedToken := generateHMACToken(data, secretKey)
	return hmac.Equal([]byte(token), []byte(expectedToken))
}

// Storage management
func (s *TokenStorage) addToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cleanup expired tokens if storage is full
	if len(s.tokens) >= maxTokens {
		s.cleanupExpiredTokens()
	}

	// If still full, return error
	if len(s.tokens) >= maxTokens {
		return ErrStorageFull
	}

	s.tokens[token] = tokenInfo{
		createdAt: time.Now(),
	}

	return nil
}

func (s *TokenStorage) cleanupExpiredTokens() {
	now := time.Now()
	for token, info := range s.tokens {
		if now.Sub(info.createdAt) > s.tokenLifetime {
			delete(s.tokens, token)
		}
	}
}

func (s *TokenStorage) verifyAndRemoveToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	info, exists := s.tokens[token]
	if !exists {
		return ErrTokenNotFound
	}

	if time.Since(info.createdAt) > s.tokenLifetime {
		delete(s.tokens, token)
		return ErrTokenExpired
	}

	// Remove token after successful verification
	delete(s.tokens, token)
	return nil
}

// Public API
func GenerateToken(data, secretKey string) (string, error) {
	token := generateHMACToken(data, secretKey)
	if err := storage.addToken(token); err != nil {
		return "", err
	}
	return token, nil
}

func VerifyToken(data, token, secretKey string) error {
	if !verifyHMACToken(data, token, secretKey) {
		return ErrInvalidToken
	}

	return storage.verifyAndRemoveToken(token)
}
