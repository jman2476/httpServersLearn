package auth

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrMissingAPIKey   = errors.New("missing API key")
	ErrMalformedAPIKey = errors.New("malformed API key")
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKeystr := headers.Get("Authorization")
	if apiKeystr == "" {
		return "", ErrMissingAPIKey
	}

	key, ok := strings.CutPrefix(apiKeystr, "ApiKey ")
	if !ok {
		return "", ErrMalformedAPIKey
	}

	return key, nil
}
