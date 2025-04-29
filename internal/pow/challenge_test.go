package pow

import (
	"testing"
)

func TestVerifySolution(t *testing.T) {
	tests := []struct {
		name       string
		challenge  string
		solution   string
		difficulty int
		want       bool
	}{
		{
			name:       "valid solution 1",
			challenge:  "16ed1c0d9a938dfe45e9c4feff8b44b6e6d42b911de47789a5c6ad8d707812cc",
			solution:   "424117",
			difficulty: 33,
			want:       true,
		},
		{
			name:       "valid solution 2",
			challenge:  "513aea0b1a195d852da6838860e251c7c00843d89a0367cc9ad8c55f37da88e0",
			solution:   "199903",
			difficulty: 12,
			want:       true,
		},
		{
			name:       "invalid solution 1",
			challenge:  "6c913a09dbb97b7091f923z04ff226x2983543584b10d65aa8b47d2e2c5v5d6c",
			solution:   "735994",
			difficulty: 11,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := VerifySolution(tt.challenge, tt.solution, tt.difficulty)
			if got != tt.want {
				t.Errorf("VerifySolution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSolveChallenge(t *testing.T) {
	challenge, err := GenerateChallenge(30)
	if err != nil {
		t.Fatalf("GenerateChallenge() failed: %v", err)
	}

	solution, err := SolveChallenge(challenge.Challenge, challenge.Difficulty)
	if err != nil {
		t.Fatalf("SolveChallenge() failed: %v", err)
	}

	if VerifySolution(challenge.Challenge, solution, challenge.Difficulty) {
		t.Logf("Found solution: %v for challenge: %v", solution, challenge.Challenge)
	} else {
		t.Errorf("VerifySolution() failed for found solution: %v", solution)
	}
}
