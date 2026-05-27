package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	// Generate 32 bytes (256 bits) of random data
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert to hex string
	tokenString := hex.EncodeToString(randomBytes)
	return tokenString, nil
}
