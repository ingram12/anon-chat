package token

import (
	"testing"
)

func TestRandomKey(t *testing.T) {
	key1 := RandomKey()

	if len(key1) != 36 { // Because hex.EncodeToString doubles the length
		t.Errorf("Expected key length to be %d, got %d", 36, len(key1))
	}

	key2 := RandomKey()

	if key1 == key2 {
		t.Error("Expected keys to be different, but they were the same")
	}
}
