package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateUserParams struct {
	FullName        string
	Email           string
	TelephoneNumber string
	Password        string
	UserRole        database.Role
}
type UpdateUserParams struct {
	FullName        string
	Email           string
	TelephoneNumber string
	UserId          int64
}
type UserRepository interface {
	Create(ctx context.Context, params CreateUserParams) (database.User, error)
	GetByEmail(ctx context.Context, email string) (database.User, error)
	GetById(ctx context.Context, id int64) (database.User, error)
	Update(ctx context.Context, params UpdateUserParams) error
	UpdateProfilePicture(ctx context.Context, profilePictureURL string, userId int64) error
}

type userRepository struct {
	store *database.Store
}

func NewUserRepository(store *database.Store) UserRepository {
	return &userRepository{
		store,
	}
}

func (r *userRepository) Create(ctx context.Context, params CreateUserParams) (database.User, error) {
	return r.store.CreateUser(ctx, database.CreateUserParams{
		FullName:        params.FullName,
		Email:           params.Email,
		TelephoneNumber: params.TelephoneNumber,
		Password:        params.Password,
		UserRole:        params.UserRole,
	})
}

func (r *userRepository) Update(ctx context.Context, params UpdateUserParams) error {
	return r.store.UpdateUser(ctx, database.UpdateUserParams{
		FullName:        params.FullName,
		Email:           params.Email,
		TelephoneNumber: params.TelephoneNumber,
		UserID:          params.UserId,
	})
}

func (r *userRepository) UpdateProfilePicture(ctx context.Context, profilePictureURL string, userId int64) error {
	return r.store.UpdateProfilePicture(ctx, database.UpdateProfilePictureParams{
		ProfileImageUrl: profilePictureURL,
		UserID:          userId,
	})
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return r.store.GetUserByEmail(ctx, email)
}

func (r *userRepository) GetById(ctx context.Context, userId int64) (database.User, error) {
	return r.store.GetUserById(ctx, userId)
}
