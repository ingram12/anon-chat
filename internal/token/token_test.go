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
		time     string
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
			got := generateHMACToken(tt.input, testSecretKey, tt.time)
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
		time     string
		expected bool
	}{
		{
			name:     "valid token",
			data:     "test",
			token:    "6e506aa16b7f8f80a8ea463b47e604ee07c8f362100bb25d3550f9319a27ccff",
			time:     "",
			expected: true,
		},
		{
			name:     "invalid token",
			data:     "test",
			token:    "invalidtoken",
			time:     "",
			expected: false,
		},
		{
			name:     "empty data",
			data:     "",
			time:     "",
			token:    "ee4c5a4bf309bd5c06cfb753ac303880aa5f4611cd97d35300e1bef35b2b539e",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verifyHMACToken(tt.data, tt.token, testSecretKey, tt.time)
			if got != tt.expected {
				t.Errorf("verifyHMACToken() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTokenStorage(t *testing.T) {
	// Override the global storage for testing purposes.
	originalStorage := storage
	storage = &Storage{
		tokenLifetime: 1 * time.Second, // Short lifetime for testing
	}
	defer func() {
		storage = originalStorage // Restore the original storage after the test.
	}()

	data := "test_data"
	// Generate a token and verify it.
	token1, err := GenerateToken(data, testSecretKey)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	err = VerifyToken(data, token1, testSecretKey)
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}

	// Wait for a short time and generate another token.
	time.Sleep(400 * time.Millisecond)
	token2, err := GenerateToken(data, testSecretKey)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// The token should be the same as the first one.
	if token1 != token2 {
		t.Errorf("GenerateToken() = %v, want %v", token2, token1)
	}

	// Wait for the token to expire.
	time.Sleep(1 * time.Second)
	token3, err := GenerateToken(data, testSecretKey)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	// The token should be different from the first one.
	if token1 == token3 {
		t.Errorf("GenerateToken() = %v, want a different token (expired)", token3)
	}

	// Verify the new token.
	err = VerifyToken(data, token3, testSecretKey)
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}

	// Test invalid token
	err = VerifyToken(data, "invalid_token", testSecretKey)
	if err == nil {
		t.Fatalf("VerifyToken() should return error, but got nil")
	}
}
