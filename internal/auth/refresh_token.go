package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	refresh := make([]byte, 32)
	rand.Read(refresh)
	return hex.EncodeToString(refresh)
}
