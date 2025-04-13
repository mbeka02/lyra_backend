package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"time"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/objstore"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
	"github.com/mbeka02/lyra_backend/internal/streamsdk"
)

const maxFileSize = 1024 * 1024 * 10 // 10 MB Limit
var allowedImageTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/bmp":     true,
	"image/tiff":    true,
	"image/svg+xml": true,
	"image/x-icon":  true, // For .ico files
	"image/heif":    true,
	"image/heic":    true,
}

type UserService interface {
	CreateUser(ctx context.Context, req model.CreateUserRequest) (model.AuthResponse, error)
	GetUser(ctx context.Context, userId int64) (model.UserResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error)
	UpdateUser(ctx context.Context, req model.UpdateUserRequest, userId int64) error
	UpdateProfilePicture(ctx context.Context, fileHeader *multipart.FileHeader, userId int64) error
}

type userService struct {
	userRepo            repository.UserRepository
	authMaker           auth.Maker
	streamClient        *streamsdk.StreamClient
	imgStorage          objstore.Storage
	accessTokenDuration time.Duration
}

// NewUserService initializes a new UserService.
func NewUserService(
	userRepo repository.UserRepository,
	authMaker auth.Maker,
	streamClient *streamsdk.StreamClient,
	imgStorage objstore.Storage,
	accessTokenDuration time.Duration,
) UserService {
	return &userService{
		userRepo:            userRepo,
		authMaker:           authMaker,
		streamClient:        streamClient,
		imgStorage:          imgStorage,
		accessTokenDuration: accessTokenDuration,
	}
}

// TODO:make the return value a ptr
func (s *userService) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.AuthResponse, error) {
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to process password:%v", err)
	}

	user, err := s.userRepo.Create(ctx, repository.CreateUserParams{
		FullName:        req.Fullname,
		Email:           req.Email,
		Password:        passwordHash,
		UserRole:        req.Role,
		TelephoneNumber: req.TelephoneNumber,
		DateOfBirth:     req.DateOfBirth,
	})
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to create user:%v", err)
	}

	userResponse := model.NewUserResponse(user)
	// auth token
	accessToken, err := s.authMaker.Create(user.UserID, user.Email, string(user.UserRole), s.accessTokenDuration)
	if err != nil {
		return model.AuthResponse{}, err
	}
	// getstream
	err = s.streamClient.CreateUser(ctx, streamsdk.CreateStreamUserParams{
		UserID: user.UserID,
		Name:   user.FullName,
		Email:  user.Email,
	})
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("stream client error : %v", err)
	}
	getStreamToken, err := s.streamClient.CreateToken(fmt.Sprintf("%d", user.UserID))
	if err != nil {
		return model.AuthResponse{}, err
	}
	return model.AuthResponse{
		AccessToken:    accessToken,
		GetStreamToken: getStreamToken,
		User:           userResponse,
	}, nil
}

func (s *userService) GetUser(ctx context.Context, userId int64) (model.UserResponse, error) {
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return model.UserResponse{}, fmt.Errorf("unable to get user details:%v", err)
	}
	return model.NewUserResponse(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, req model.UpdateUserRequest, userId int64) error {
	return s.userRepo.Update(ctx, repository.UpdateUserParams{
		Email:           req.Email,
		TelephoneNumber: req.TelephoneNumber,
		FullName:        req.FullName,
		UserId:          userId,
	})
}

func (s *userService) UpdateProfilePicture(ctx context.Context, fileHeader *multipart.FileHeader, userId int64) error {
	err := isValidImageFile(fileHeader)
	if err != nil {
		return err
	}
	user, err := s.userRepo.GetById(ctx, userId)
	if err != nil {
		return fmt.Errorf("unable to get user details:%v", err)
	}
	objectName, err := objNameFromURL(user.ProfileImageUrl, fileHeader.Filename)
	if err != nil {
		return err
	}
	imageURL, err := s.imgStorage.Upload(ctx, objectName, fileHeader)
	if err != nil {
		return fmt.Errorf("unable to upload the image:%v", err)
	}
	return s.userRepo.UpdateProfilePicture(ctx, imageURL, userId)
}

func (s *userService) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return model.AuthResponse{}, errors.New("account does not exist")
	}

	if err := auth.ComparePassword(req.Password, user.Password); err != nil {
		return model.AuthResponse{}, errors.New("the password is invalid")
	}

	userResponse := model.NewUserResponse(user)
	accessToken, err := s.authMaker.Create(user.UserID, user.Email, string(user.UserRole), s.accessTokenDuration)
	if err != nil {
		return model.AuthResponse{}, err
	}
	getStreamToken, err := s.streamClient.CreateToken(fmt.Sprintf("%d", user.UserID))
	if err != nil {
		return model.AuthResponse{}, err
	}
	return model.AuthResponse{
		AccessToken:    accessToken,
		GetStreamToken: getStreamToken,
		User:           userResponse,
	}, nil
}

func objNameFromURL(imageURL string, fileName string) (string, error) {
	// if user doesn't have imageURL - create one
	// otherwise, extract last part of URL to get cloud storage object name
	if imageURL == "" {
		objectName := fmt.Sprintf("%s_%d", fileName, time.Now().UnixNano())

		return objectName, nil
	}

	// split off last part of URL, which is the image's storage object ID
	urlPath, err := url.Parse(imageURL)
	if err != nil {
		return "", errors.New("Failed to parse objectName from imageURL")
	}

	// get "path" of url (everything after domain)
	// then get "base", the last part

	return path.Base(urlPath.Path), nil
}

func isValidImageFile(fileHeader *multipart.FileHeader) error {
	fileContentType := fileHeader.Header.Get("Content-Type")
	if fileHeader.Size > maxFileSize {
		return fmt.Errorf("the image size is too large: %v", fileHeader.Size)
	}
	if _, ok := allowedImageTypes[fileContentType]; !ok {
		return fmt.Errorf("this file format is not supported: %v", fileContentType)
	}
	return nil
}

// Convert a string to sql.NullString
func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
