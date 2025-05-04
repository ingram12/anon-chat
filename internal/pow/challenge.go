package pow

import (
	"anon-chat/internal/token"
	"crypto/sha256"
	"errors"
	"fmt"
)

func VerifyChallengeNonce(challenge, nonce string, difficulty int) bool {
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
		if VerifyChallengeNonce(challenge, nonce, difficulty) {
			return nonce, nil
		}
	}
	return "", errors.New("could not find valid nonce within maxAttempts")
}

func GenerateChallenge() string {
	return token.RandomKey()
}
