package auth

import (
	"fmt"
	"runtime"

	"github.com/alexedwards/argon2id"
)

var hashParameters = argon2id.Params{
	Memory:      4 * 1024 * 1024,
	Iterations:  1,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   32,
}

func HashPassword(password string) (string, error) {

	hash, err := argon2id.CreateHash(password, &hashParameters)
	if err != nil {
		return "", fmt.Errorf("Error hashing password: %w", err)
	}

	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
