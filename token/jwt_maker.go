package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTMaker is a JSON web token maker
type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

// NewJWTMaker creates a new JWTMaker. Maker is an inteface defined in token.maker.go
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least 32 characters")
	}

	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	// JWT create token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	Keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, Keyfunc)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// convert JWT claims into the payload
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
