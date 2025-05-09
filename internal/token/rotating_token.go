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

func NewRotatingToken(lifeTime time.Duration) *RotatingToken {
	return &RotatingToken{
		lifetime: lifeTime,
	}
}

func (s *RotatingToken) GetRotatingToken() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if the token is expired or not set
	if s.value == "" || s.expiresAt.Before(time.Now()) {
		// Generate a new token
		s.expiresAt = time.Now().Add(s.lifetime)
		s.value = RandomKey()
	}

	return s.value
}
