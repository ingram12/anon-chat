package token

import (
	"testing"
)

func TestGenerateUserToken(t *testing.T) {
	token, err := GenerateUserToken()
	if err != nil {
		t.Fatalf("Error generating user token: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	token2, err := GenerateUserToken()
	if err != nil {
		t.Fatalf("Error generating user token: %v", err)
	}

	if token2 == "" {
		t.Error("Token2 should not be empty")
	}

	if token == token2 {
		t.Error("Token should be different on each call")
	}
}
