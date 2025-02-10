package model

import (
	"database/sql"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateUserRequest struct {
	Fullname        string        `json:"full_name" validate:"required,min=2"`
	Email           string        `json:"email" validate:"required,email"`
	TelephoneNumber string        `json:"telephone_number" validate:"required,max=15"`
	Password        string        `json:"password" validate:"required,min=8"`
	Role            database.Role `json:"role" validate:"required"`
}
type UpdateUserRequest struct {
	FullName        string `json:"full_name" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	TelephoneNumber string `json:"telephone_number" validate:"required,max=15"`
}
type UserResponse struct {
	Fullname        string         `json:"full_name" validate:"required"`
	Email           string         `json:"email" validate:"required,email"`
	TelephoneNumber string         `json:"telephone_number" validate:"required,max=15"`
	Password        string         `json:"password" validate:"required,min=8"`
	Role            database.Role  `json:"role" validate:"required"`
	ProfileImageURL sql.NullString `json:"profile_image_url"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func NewUserResponse(user database.User) UserResponse {
	return UserResponse{
		Fullname:        user.FullName,
		Email:           user.Email,
		TelephoneNumber: user.TelephoneNumber,
		Role:            user.UserRole,
		ProfileImageURL: user.ProfileImageUrl,
	}
}
