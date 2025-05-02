package token

import (
	"testing"
)

func TestGenerateAndVerifyHMACToken(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		secretKey string
	}{
		{
			name:      "simple token",
			data:      "test-data",
			secretKey: "secret-key",
		},
		{
			name:      "empty data",
			data:      "",
			secretKey: "secret-key",
		},
		{
			name:      "with special characters",
			data:      "test|data|with|pipes",
			secretKey: "secret-key!@#$%^&*()",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate token
			token := GenerateHMACToken(tt.data, tt.secretKey)
			if token == "" {
				t.Error("GenerateHMACToken() returned empty token")
			}

			// Verify valid token
			if !VerifyHMACToken(tt.data, token, tt.secretKey) {
				t.Error("VerifyHMACToken() returned false for valid token")
			}

			// Verify with modified data
			if VerifyHMACToken(tt.data+"modified", token, tt.secretKey) {
				t.Error("VerifyHMACToken() returned true for modified data")
			}

			// Verify with modified token
			if VerifyHMACToken(tt.data, token+"modified", tt.secretKey) {
				t.Error("VerifyHMACToken() returned true for modified token")
			}

			// Verify with different secret key
			if VerifyHMACToken(tt.data, token, tt.secretKey+"different") {
				t.Error("VerifyHMACToken() returned true for different secret key")
			}
		})
	}
}

func TestHMACTokenDeterministic(t *testing.T) {
	data := "test-data"
	secretKey := "secret-key"

	// Generate two tokens with same input
	token1 := GenerateHMACToken(data, secretKey)
	token2 := GenerateHMACToken(data, secretKey)

	if token1 != token2 {
		t.Error("GenerateHMACToken() should be deterministic for same input")
	}
}
