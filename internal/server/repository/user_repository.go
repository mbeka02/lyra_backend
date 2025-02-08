package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateUserParams struct {
	FullName string
	Email    string
	Password string
}
type UserRepository interface {
	Create(ctx context.Context, params CreateUserParams) (database.User, error)
	GetByEmail(ctx context.Context, email string) (database.User, error)
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
		FullName: params.FullName,
		Email:    params.Email,
		Password: params.Password,
	})
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return r.store.GetUserByEmail(ctx, email)
}
