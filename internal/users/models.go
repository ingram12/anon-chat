package users

// GetFirstChallengeResponse represents the response for getting a new challenge
type GetFirstChallengeResponse struct {
	Challenge  string `json:"challenge"`
	Token      string `json:"token"`
	Difficulty int    `json:"difficulty"`
}

// SolveFirstChallengeRequest represents the request for solving a challenge
type SolveFirstChallengeRequest struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Solution  string `json:"solution"`
}

// SolveFirstChallengeResponse represents the response for solving a challenge
type SolveFirstChallengeResponse struct {
	UserID           string `json:"user_id"`
	CurrentChallenge string `json:"current_challenge"`
}

// RegisterUserRequest represents the request for user registration
type RegisterUserRequest struct {
	UserID    string   `json:"user_id"`
	Nickname  string   `json:"nickname"`
	Tags      []string `json:"tags"`
	PublicKey string   `json:"public_key"`
	Challenge string   `json:"challenge"`
	Solution  string   `json:"solution"`
}

// RegisterUserResponse represents the response for user registration
type RegisterUserResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
