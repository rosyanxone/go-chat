package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GeneratePlainToken() (string, error) {
	b := make([]byte, 32) // 32 bytes = 256-bit random
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func BuildPlainTextToken(tokenID uint64, plainToken string) string {
	return fmt.Sprintf("%d|%s", tokenID, plainToken)
}
