package token

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	ChallengeLength = 16 // Length of the challenge in bytes
)

func RandomKey() (string, error) {
	b := make([]byte, ChallengeLength)
	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
