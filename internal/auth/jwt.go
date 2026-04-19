package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var malformedHeaderErr = errors.New("Malformed token in header")
var missingBearerTokenErr = errors.New("Header does not contain token")

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy-access",
	}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing token string %w", err)
	}

	idString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error getting claims subject %w", err)
	}

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Error parsing uuid string %w", err)
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearerStr := headers.Get("Authorization")
	if bearerStr == "" {
		return "", missingBearerTokenErr
	}

	token, ok := strings.CutPrefix(bearerStr, "Bearer ")
	if !ok {
		return "", malformedHeaderErr
	}

	return token, nil
}
