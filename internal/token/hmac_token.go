package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHMACToken(data, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func VerifyHMACToken(data, token, secretKey string) bool {
	expectedToken := GenerateHMACToken(data, secretKey)
	return hmac.Equal([]byte(token), []byte(expectedToken))
}
