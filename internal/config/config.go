package config

import "os"

type Config struct {
	TokenSecretKey string
	RandomKeySize  int
	SecondLifeTime int
}

func NewConfig() *Config {
	secretKey := os.Getenv("TOKEN_SECRET_KEY")

	if secretKey == "" {
		secretKey = "DEFAULT-SECRET-KEY-ALARMfgsjffsr"
	}

	return &Config{
		TokenSecretKey: secretKey,
		RandomKeySize:  16,
		SecondLifeTime: 120,
	}
}
