package config

import (
	"os"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	// Clear any existing TOKEN_SECRET_KEY environment variable
	os.Unsetenv("TOKEN_SECRET_KEY")

	tests := []struct {
		name   string
		isDev  bool
		envKey string
		want   *Config
	}{
		{
			name:  "Development mode with no env key",
			isDev: true,
			want: &Config{
				TokenSecretKey:           "DEFAULT-SECRET-KEY-ALARMfgsjffsr",
				FirstChallengeDifficulty: 1,
				RotatingTokenLifeTime:    180 * time.Second,
				UserInactivityTimeout:    1800 * time.Second,
			},
		},
		{
			name:   "Development mode with custom env key",
			isDev:  true,
			envKey: "custom-secret-key",
			want: &Config{
				TokenSecretKey:           "custom-secret-key",
				FirstChallengeDifficulty: 1,
				RotatingTokenLifeTime:    180 * time.Second,
				UserInactivityTimeout:    1800 * time.Second,
			},
		},
		{
			name:  "Production mode",
			isDev: false,
			want: &Config{
				FirstChallengeDifficulty: 9999,
				RotatingTokenLifeTime:    180 * time.Second,
				UserInactivityTimeout:    1800 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" {
				os.Setenv("TOKEN_SECRET_KEY", tt.envKey)
				defer os.Unsetenv("TOKEN_SECRET_KEY")
			} else {
				os.Unsetenv("TOKEN_SECRET_KEY")
			}

			got := NewConfig(tt.isDev)

			// In production mode, we generate a random key, so we can't test the exact value
			if !tt.isDev && tt.envKey == "" {
				if got.TokenSecretKey == "" {
					t.Error("Expected random token key in production, got empty string")
				}
				// Copy the generated key to our want struct for comparison
				tt.want.TokenSecretKey = got.TokenSecretKey
			}

			if got.TokenSecretKey != tt.want.TokenSecretKey {
				t.Errorf("TokenSecretKey = %v, want %v", got.TokenSecretKey, tt.want.TokenSecretKey)
			}
			if got.FirstChallengeDifficulty != tt.want.FirstChallengeDifficulty {
				t.Errorf("FirstChallengeDifficulty = %v, want %v", got.FirstChallengeDifficulty, tt.want.FirstChallengeDifficulty)
			}
			if got.RotatingTokenLifeTime != tt.want.RotatingTokenLifeTime {
				t.Errorf("RotatingTokenLifeTime = %v, want %v", got.RotatingTokenLifeTime, tt.want.RotatingTokenLifeTime)
			}
			if got.UserInactivityTimeout != tt.want.UserInactivityTimeout {
				t.Errorf("UserInactivityTimeout = %v, want %v", got.UserInactivityTimeout, tt.want.UserInactivityTimeout)
			}
		})
	}
}
