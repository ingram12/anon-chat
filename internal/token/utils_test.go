package token

import (
	"testing"
)

func TestRandomKey(t *testing.T) {
	key1, err := RandomKey()
	if err != nil {
		t.Fatalf("Error generating random key: %v", err)
	}

	if len(key1) != ChallengeLength*2 { // Because hex.EncodeToString doubles the length
		t.Errorf("Expected key length to be %d, got %d", ChallengeLength*2, len(key1))
	}

	key2, err := RandomKey()
	if err != nil {
		t.Fatalf("Error generating random key: %v", err)
	}

	if key1 == key2 {
		t.Error("Expected keys to be different, but they were the same")
	}
}
