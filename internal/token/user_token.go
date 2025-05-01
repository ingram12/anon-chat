package token

func GenerateUserToken() (string, error) {
	return RandomKey()
}
