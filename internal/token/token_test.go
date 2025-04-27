package token

import (
	"testing"
	"time"
)

const testSecretKey = "fgsjffsrujJJHJHGOBJWHQP'[]KKK"

func TestGenerateHMACToken(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "ee4c5a4bf309bd5c06cfb753ac303880aa5f4611cd97d35300e1bef35b2b539e",
		},
		{
			name:     "simple string",
			input:    "test",
			expected: "6e506aa16b7f8f80a8ea463b47e604ee07c8f362100bb25d3550f9319a27ccff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateHMACToken(tt.input, testSecretKey)
			if got != tt.expected {
				t.Errorf("generateHMACToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestVerifyHMACToken(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		token    string
		expected bool
	}{
		{
			name:     "valid token",
			data:     "test",
			token:    "6e506aa16b7f8f80a8ea463b47e604ee07c8f362100bb25d3550f9319a27ccff",
			expected: true,
		},
		{
			name:     "invalid token",
			data:     "test",
			token:    "invalidtoken",
			expected: false,
		},
		{
			name:     "empty data",
			data:     "",
			token:    "ee4c5a4bf309bd5c06cfb753ac303880aa5f4611cd97d35300e1bef35b2b539e",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verifyHMACToken(tt.data, tt.token, testSecretKey)
			if got != tt.expected {
				t.Errorf("verifyHMACToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTokenStorage(t *testing.T) {
	storage := &TokenStorage{
		tokens:        make(map[string]tokenInfo),
		tokenLifetime: 1 * time.Second, // Short lifetime for testing
	}

	t.Run("add and verify token", func(t *testing.T) {
		token := "test_token"
		if err := storage.addToken(token); err != nil {
			t.Fatalf("addToken() error = %v", err)
		}

		if err := storage.verifyAndRemoveToken(token); err != nil {
			t.Errorf("verifyAndRemoveToken() error = %v", err)
		}

		// Token should be removed after verification
		if err := storage.verifyAndRemoveToken(token); err != ErrTokenNotFound {
			t.Errorf("verifyAndRemoveToken() error = %v, want %v", err, ErrTokenNotFound)
		}
	})

	t.Run("expired token", func(t *testing.T) {
		token := "expired_token"
		storage.tokens[token] = tokenInfo{
			createdAt: time.Now().Add(-storage.tokenLifetime - time.Second),
		}

		if err := storage.verifyAndRemoveToken(token); err != ErrTokenExpired {
			t.Errorf("verifyAndRemoveToken() error = %v, want %v", err, ErrTokenExpired)
		}
	})

	t.Run("storage full", func(t *testing.T) {
		// Fill storage
		for i := 0; i < maxTokens; i++ {
			token := string(rune('a' + i))
			storage.tokens[token] = tokenInfo{
				createdAt: time.Now(),
			}
		}

		// Try to add one more token
		if err := storage.addToken("extra_token"); err != ErrStorageFull {
			t.Errorf("addToken() error = %v, want %v", err, ErrStorageFull)
		}
	})
}

func TestPublicAPI(t *testing.T) {
	// Create a new storage with short lifetime for testing
	testStorage := &TokenStorage{
		tokens:        make(map[string]tokenInfo),
		tokenLifetime: 1 * time.Second,
	}
	// Replace global storage for testing
	oldStorage := storage
	storage = testStorage
	defer func() { storage = oldStorage }()

	t.Run("generate and verify token", func(t *testing.T) {
		data := "test_data"
		token, err := GenerateToken(data, testSecretKey)
		if err != nil {
			t.Fatalf("GenerateToken() error = %v", err)
		}

		if err := VerifyToken(data, token, testSecretKey); err != nil {
			t.Errorf("VerifyToken() error = %v", err)
		}
	})

	t.Run("verify invalid token", func(t *testing.T) {
		data := "test_data"
		token := "invalid_token"
		if err := VerifyToken(data, token, testSecretKey); err != ErrInvalidToken {
			t.Errorf("VerifyToken() error = %v, want %v", err, ErrInvalidToken)
		}
	})

	t.Run("verify expired token", func(t *testing.T) {
		data := "test_data"
		token, err := GenerateToken(data, testSecretKey)
		if err != nil {
			t.Fatalf("GenerateToken() error = %v", err)
		}

		// Wait for token to expire
		time.Sleep(testStorage.tokenLifetime + time.Second)

		if err := VerifyToken(data, token, testSecretKey); err != ErrTokenExpired {
			t.Errorf("VerifyToken() error = %v, want %v", err, ErrTokenExpired)
		}
	})
}
