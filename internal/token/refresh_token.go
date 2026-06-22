package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateRefreshToken() string {
	return rand.Text()
}

func HashRefreshToken(refreshToken string) string {
	tokenHash := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(tokenHash[:])
}
