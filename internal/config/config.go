package config

import (
	"os"
	"time"
)

type Config struct {
	TokenSecretKey           string
	FirstChallengeDifficulty int
	RotatingTokenLifeTime    time.Duration
	UserInactivityTimeout    time.Duration
}

func NewConfig() *Config {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")

	if secretKey == "" {
		secretKey = "DEFAULT-SECRET-KEY-ALARMfgsjffsr"
	}

	return &Config{
		TokenSecretKey:           secretKey,
		FirstChallengeDifficulty: 300,
		RotatingTokenLifeTime:    120 * time.Second,
		UserInactivityTimeout:    1800 * time.Second,
	}
}
