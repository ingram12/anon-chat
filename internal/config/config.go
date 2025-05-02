package config

import "os"

type Config struct {
	TokenSecretKey           string
	FirstChallengeDifficulty int
	SecondLifeTime           int
	UserInactivityTimeout    int // in seconds
}

func NewConfig() *Config {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")

	if secretKey == "" {
		secretKey = "DEFAULT-SECRET-KEY-ALARMfgsjffsr"
	}

	return &Config{
		TokenSecretKey:           secretKey,
		FirstChallengeDifficulty: 300,
		SecondLifeTime:           120,
		UserInactivityTimeout:    1800,
	}
}
