package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	Scope     []string  `json:"scope"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayLoad(username string, role string, scope []string, duration time.Duration) (*Payload, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        token,
		Username:  username,
		Role:      role,
		Scope:     scope,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}

var ErrExpiredToken = errors.New("token has expired")

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
