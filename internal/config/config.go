package config

import (
	"anon-chat/internal/token"
	"os"
	"time"
)

type Config struct {
	TokenSecretKey           string
	FirstChallengeDifficulty int
	RotatingTokenLifeTime    time.Duration
	UserInactivityTimeout    time.Duration
}

func NewConfig(isdev bool) *Config {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")

	if secretKey == "" {
		if isdev {
			secretKey = "DEFAULT-SECRET-KEY-ALARMfgsjffsr"
		} else {
			secretKey = token.RandomKey()
		}
	}

	firstChallengeDifficulty := 9999
	if isdev {
		firstChallengeDifficulty = 1
	}

	return &Config{
		TokenSecretKey:           secretKey,
		FirstChallengeDifficulty: firstChallengeDifficulty,
		RotatingTokenLifeTime:    180 * time.Second,
		UserInactivityTimeout:    180 * time.Second,
	}
}
