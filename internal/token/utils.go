package token

import "github.com/google/uuid"

// RandomKey generates a random key using UUID v4
func RandomKey() string {
	return uuid.New().String()
}
