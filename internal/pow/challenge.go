package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
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

func VerifySolution(challenge, nonce string, difficulty int) bool {
	hash := []byte(challenge + nonce)

	for i := 0; i < difficulty; i++ {
		sum := sha256.Sum256(hash)
		hash = sum[:]
	}

	return hash[0] == 0
}

func SolveChallenge(challenge string, difficulty int) (string, error) {
	for i := range 100000000 {
		nonce := fmt.Sprintf("%d", i)
		if VerifySolution(challenge, nonce, difficulty) {
			return nonce, nil
		}
	}
	return "", errors.New("could not find valid nonce within maxAttempts")
}
