package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func GenRandKey(length int32) (string, error) {
	byts := make([]byte, length)
	if _, err := rand.Read(byts); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(byts), nil
}
