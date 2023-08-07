package authentication

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomState() (string, error) {
	const stateBytes = 32

	// Generate random bytes
	randomBytes := make([]byte, stateBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to a base64 string
	state := base64.URLEncoding.EncodeToString(randomBytes)

	return state, nil
}
