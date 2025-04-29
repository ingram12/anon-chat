package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	ChallengeLength = 32
)

type Challenge struct {
	Challenge  string
	Difficulty int
}

func GenerateChallenge(difficulty int) (Challenge, error) {
	b := make([]byte, ChallengeLength)
	if _, err := rand.Read(b); err != nil {
		return Challenge{}, err
	}

	return Challenge{
		Challenge:  hex.EncodeToString(b),
		Difficulty: difficulty,
	}, nil
}

func VerifySolution(challenge, nonce string, difficulty int) bool {
	hash := []byte(challenge + nonce)

	for i := 0; i < difficulty; i++ {
		sum := sha256.Sum256(hash)
		hash = sum[:]
	}

	return hash[0] == 0 && hash[1] == 0
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
