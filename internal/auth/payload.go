package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    int64     `json:"user_id"`
	Role      string    `json:"role"`
	jwt.RegisteredClaims
}

func NewPayload(userId int64, email, role string, duration time.Duration) *Payload {
	return &Payload{
		UserID:    userId,
		Email:     email,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}
}
