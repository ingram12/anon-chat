package token

import (
	"testing"
)

func TestGenerateUserToken(t *testing.T) {
	token := GenerateUserToken()

	if token == "" {
		t.Error("Token should not be empty")
	}

	token2 := GenerateUserToken()

	if token2 == "" {
		t.Error("Token2 should not be empty")
	}

	if token == token2 {
		t.Error("Token should be different on each call")
	}
}
