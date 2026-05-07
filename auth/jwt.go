package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken  = errors.New("token has expired")
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// jwtSecret is the default signing secret. Override at startup via SetSecret.
var jwtSecret = []byte("s-ui-secret-key")

// SetSecret allows overriding the default JWT secret (e.g. from config).
// An empty string is ignored so the existing secret is preserved.
func SetSecret(secret string) {
	if secret != "" {
		jwtSecret = []byte(secret)
	}
}

// GenerateToken creates a signed JWT for the given username with the given TTL.
func GenerateToken(username string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			// NotBefore ensures the token cannot be used before it is issued.
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken validates the token string and returns the embedded Claims.
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
