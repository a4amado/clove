package apiguard

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func randKey() (string, error) {
	byts := make([]byte, 64)
	_, err := rand.Read(byts)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(byts), nil

}

func RandomSecretKey() (string, error) {
	key, err := randKey()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("sk-%s", key), nil
}
