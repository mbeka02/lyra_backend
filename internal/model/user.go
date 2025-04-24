package model

import (
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateUserRequest struct {
	Fullname        string        `json:"full_name" validate:"required,min=2"`
	Email           string        `json:"email" validate:"required,email"`
	TelephoneNumber string        `json:"telephone_number" validate:"required,max=15"`
	Password        string        `json:"password" validate:"required,min=8"`
	Role            database.Role `json:"role" validate:"required"`
	DateOfBirth     time.Time     `json:"date_of_birth" validate:"required"`
}
type UpdateUserRequest struct {
	FullName        string `json:"full_name" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	TelephoneNumber string `json:"telephone_number" validate:"required,max=15"`
}
type UserResponse struct {
	UserId          int64         `json:"user_id"`
	Fullname        string        `json:"full_name"`
	Email           string        `json:"email"`
	TelephoneNumber string        `json:"telephone_number" `
	Role            database.Role `json:"role" `
	ProfileImageURL string        `json:"profile_image_url"`
	IsOnboarded     bool          `json:"is_onboarded"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	AccessToken    string       `json:"access_token"`
	GetStreamToken string       `json:"get_stream_token"`
	User           UserResponse `json:"user"`
}

func NewUserResponse(user *database.User) UserResponse {
	return UserResponse{
		UserId:          user.UserID,
		Fullname:        user.FullName,
		Email:           user.Email,
		TelephoneNumber: user.TelephoneNumber,
		Role:            user.UserRole,
		ProfileImageURL: user.ProfileImageUrl,
		IsOnboarded:     user.IsOnboarded,
	}
}
