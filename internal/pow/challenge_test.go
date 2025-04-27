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
			challenge:  "1e37c85282826db71ae3fcc6307639073c2a985075604d2b2d1249ee51d9f674",
			solution:   "41",
			difficulty: 10,
			want:       true,
		},
		{
			name:       "valid solution 2",
			challenge:  "d176f4ac9b0f38f706ae4628fa62f5350bd261d9344abfff04bbb49432725989",
			solution:   "497",
			difficulty: 4000,
			want:       true,
		},
		{
			name:       "invalid solution 1",
			challenge:  "6c913a09dbb97b7091f923z04ff226x2983543584b10d65aa8b47d2e2c5v5d6c",
			solution:   "7394",
			difficulty: 48,
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
