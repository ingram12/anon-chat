package token

import (
	"testing"
	"time"
)

func TestRotatingToken_GetRotatingToken(t *testing.T) {
	token := NewRotatingToken(50 * time.Millisecond) // Token expires in 1 second

	// First call should generate a new token
	token1 := token.GetRotatingToken()

	if token1 == "" {
		t.Error("Token should not be empty")
	}

	// Second call immediately after should return the same token
	token2 := token.GetRotatingToken()

	if token1 != token2 {
		t.Errorf("Expected token to be the same, got %s and %s", token1, token2)
	}

	// Wait for the token to expire
	time.Sleep(51 * time.Millisecond)

	// Third call after expiration should generate a new token
	token3 := token.GetRotatingToken()

	if token1 == token3 {
		t.Error("Expected token to be different after expiration")
	}
}

func TestNewRotatingToken(t *testing.T) {
	lifetime := 60 * time.Second // seconds
	token := NewRotatingToken(lifetime)

	if token.lifetime != lifetime {
		t.Errorf("Expected lifetime to be %v, got %v", lifetime, token.lifetime)
	}
}
