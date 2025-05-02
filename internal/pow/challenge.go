package pow

import (
	"anon-chat/internal/token"
	"anon-chat/internal/users"
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

func VerifyFullChallenge(
	challenge string,
	nonce string,
	token string,
	difficulty int,
	secretKey string,
	storage *users.UserStorage,
	rotatingToken *token.RotatingToken,
) error {
	isUserExist := storage.IsUserExist(token)
	if isUserExist {
		return errors.New("challenge already solved")
	}

	userToken := token
	globalToken, err := rotatingToken.GetRotatingToken()
	if err != nil {
		return errors.New("global token generation failed")
	}

	isVerified := VerifyFirstChallenge(userToken, globalToken, challenge, secretKey)
	if !isVerified {
		return errors.New("challenge verification failed")
	}

	if !VerifyChallengeNonce(challenge, nonce, int(difficulty)) {
		return errors.New("invalid PoW nonce")
	}

	return nil
}
