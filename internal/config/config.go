package config

import "os"

type Config struct {
	TokenSecretKey string
}

func NewConfig() *Config {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")

	if secretKey == "" {
		secretKey = "fgsjffsrujJJHJHGOBJWHQP'[]KKK"
	}

	return &Config{
		TokenSecretKey: secretKey,
	}
}
