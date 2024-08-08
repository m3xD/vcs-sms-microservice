package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

func NewJWT(secretKey string) (Maker, error) {
	return &JWTMaker{secretKey: secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, role string, scope []string, duration time.Duration) (string, error) {
	payload, err := NewPayLoad(username, role, scope, duration)

	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

var ErrInvalidToken = errors.New("token is invalid")

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		ver, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(ver.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
