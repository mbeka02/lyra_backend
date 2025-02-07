package auth

import "time"

type Maker interface {
	Create(email string, userId int64, duration time.Duration) (string, error)
	Verify(tokenString string) (*Payload, error)
}
