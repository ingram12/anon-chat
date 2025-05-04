package pow

import (
	"anon-chat/internal/token"
	"testing"
)

func TestVerifyChallengeNonce(t *testing.T) {
	tests := []struct {
		name       string
		challenge  string
		nonce      string
		difficulty int
		want       bool
	}{
		{
			name:       "valid solution 1",
			challenge:  "16ed1c0d9a938dfe45e9c4feff8b44b6e6d42b911de47789a5c6ad8d707812cc",
			nonce:      "424117",
			difficulty: 33,
			want:       true,
		},
		{
			name:       "valid solution 2",
			challenge:  "513aea0b1a195d852da6838860e251c7c00843d89a0367cc9ad8c55f37da88e0",
			nonce:      "199903",
			difficulty: 12,
			want:       true,
		},
		{
			name:       "invalid solution 1",
			challenge:  "6c913a09dbb97b7091f923z04ff226x2983543584b10d65aa8b47d2e2c5v5d6c",
			nonce:      "735994",
			difficulty: 11,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifyChallengeNonce(tt.challenge, tt.nonce, tt.difficulty)
			if got != tt.want {
				t.Errorf("VerifyChallengeNonce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolveChallenge(t *testing.T) {
	challenge := token.RandomKey()

	solution, err := SolveChallenge(challenge, 30)
	if err != nil {
		t.Fatalf("SolveChallenge() failed: %v", err)
	}

	if !VerifyChallengeNonce(challenge, solution, 30) {
		t.Errorf("VerifyChallengeNonce() failed for found solution: %v", solution)
	}
}
