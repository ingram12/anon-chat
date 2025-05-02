package pow

import "testing"

func TestGenerateAndVerifyChallenge(t *testing.T) {
	userKey := "testUser"
	globalKey := "testGlobal"
	secretKey := "testSecret"

	challenge := GenerateFirstChallenge(userKey, globalKey, secretKey)
	if challenge == "" {
		t.Error("GenerateChallenge() returned an empty string")
	}

	isValid := VerifyFirstChallenge(userKey, globalKey, challenge, secretKey)
	if !isValid {
		t.Error("VerifyChallenge() returned false for a valid challenge")
	}

	// Test with a modified challenge
	invalidChallenge := challenge + "modified"
	isInvalid := VerifyFirstChallenge(userKey, globalKey, invalidChallenge, secretKey)
	if isInvalid {
		t.Error("VerifyChallenge() returned true for an invalid challenge")
	}

	// Test with a different user key
	isInvalid = VerifyFirstChallenge("differentUser", globalKey, challenge, secretKey)
	if isInvalid {
		t.Error("VerifyChallenge() returned true for a challenge with a different user key")
	}

	// Test with a different global key
	isInvalid = VerifyFirstChallenge(userKey, "differentGlobal", challenge, secretKey)
	if isInvalid {
		t.Error("VerifyChallenge() returned true for a challenge with a different global key")
	}
}
