package auth

import "time"

type Maker interface {
	Create(userId int64, email, role string, duration time.Duration) (string, error)
	Verify(tokenString string) (*Payload, error)
}
