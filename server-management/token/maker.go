package token

import "time"

type Maker interface {
	CreateToken(username string, role string, scope []string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
