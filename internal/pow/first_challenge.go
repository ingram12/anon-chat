package pow

import (
	"anon-chat/internal/token"
)

// Generates a challenge token based on the user key and the global temporary key
func GenerateFirstChallenge(userKey, globalKey, secretKey string) string {
	return token.GenerateHMACToken(
		globalKey+"|"+userKey,
		secretKey,
	)
}

// Returns true if the token was created by the server and its lifetime has not expired
func VerifyFirstChallenge(userKey, globalKey, hmacToken, secretKey string) bool {
	return token.VerifyHMACToken(
		globalKey+"|"+userKey,
		hmacToken,
		secretKey,
	)
}
