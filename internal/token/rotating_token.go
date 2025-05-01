package token

import (
	"sync"
	"time"
)

type RotatingToken struct {
	mu        sync.RWMutex
	value     string
	expiresAt time.Time
	lifetime  time.Duration
}

func NewRotatingToken(secondLifeTime int) *RotatingToken {
	return &RotatingToken{
		lifetime: time.Duration(secondLifeTime) * time.Second,
	}
}

func (s *RotatingToken) GetRotatingToken() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the token is expired or not set
	if s.value == "" || s.expiresAt.Before(time.Now()) {
		// Generate a new RotatingToken
		token, err := RandomKey()
		if err != nil {
			return "", err
		}

		s.expiresAt = time.Now().Add(s.lifetime)
		s.value = token
	}

	return s.value, nil
}
