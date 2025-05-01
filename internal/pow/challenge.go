package pow

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

func generateHMACToken(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func verifyHMACToken(data, token, secretKey string) bool {
	expectedToken := generateHMACToken(data, secretKey)
	return hmac.Equal([]byte(token), []byte(expectedToken))
}

// Generates a challenge token based on the user key and the global temporary key
func GenerateChallenge(userKey, globalKey, secretKey string) string {
	return generateHMACToken(
		globalKey+"|"+userKey,
		secretKey,
	)
}

// Returns true if the token was created by the server and its lifetime has not expired
func VerifyChallenge(userKey, globalKey, token, secretKey string) bool {
	return verifyHMACToken(
		globalKey+"|"+userKey,
		token,
		secretKey,
	)
}
