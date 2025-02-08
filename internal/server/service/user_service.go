package service

import (
	"context"
	"errors"
	"time"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, req model.CreateUserRequest) (model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error)
}

type userService struct {
	repo                repository.UserRepository
	authMaker           auth.Maker
	accessTokenDuration time.Duration
}

func NewUserService(repo repository.UserRepository, authMaker auth.Maker, duration time.Duration) UserService {
	return &userService{
		repo:                repo,
		authMaker:           authMaker,
		accessTokenDuration: duration,
	}
}

func (s *userService) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.AuthResponse, error) {
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.AuthResponse{}, errors.New("failed to process password")
	}

	user, err := s.repo.Create(ctx, repository.CreateUserParams{
		FullName: req.Fullname,
		Email:    req.Email,
		Password: passwordHash,
	})
	if err != nil {
		return model.AuthResponse{}, errors.New("failed to create user")
	}

	userResponse := model.NewUserResponse(user)
	token, err := s.authMaker.Create(user.Email, user.UserID, s.accessTokenDuration)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		AccessToken: token,
		User:        userResponse,
	}, nil
}

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return model.AuthResponse{}, errors.New("unable to find user")
	}

	if err := auth.ComparePassword(req.Password, user.Password); err != nil {
		return model.AuthResponse{}, err
	}

	userResponse := model.NewUserResponse(user)
	token, err := s.authMaker.Create(user.Email, user.UserID, s.accessTokenDuration)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		AccessToken: token,
		User:        userResponse,
	}, nil
}
