package users

func BytesToString(id [36]byte) string {
	return string(id[:])
}

func StringToBytes(userID string) [36]byte {
	var idBytes [36]byte
	copy(idBytes[:], userID)
	return idBytes
}
