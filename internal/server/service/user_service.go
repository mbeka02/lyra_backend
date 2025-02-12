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
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/objstore"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
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
	CreatePatient(ctx context.Context, req model.CreatePatientRequest, userId int64) (database.Patient, error)
	CreateSpecialist(ctx context.Context, req model.CreateSpecialistRequest, userId int64) (database.Specialist, error)
	GetUser(ctx context.Context, userId int64) (model.UserResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error)
	UpdateUser(ctx context.Context, req model.UpdateUserRequest, userId int64) error
	UpdateProfilePicture(ctx context.Context, fileHeader *multipart.FileHeader, userId int64) error
}

type userService struct {
	userRepo            repository.UserRepository
	patientRepo         repository.PatientRepository
	specialistRepo      repository.SpecialistRepository
	authMaker           auth.Maker
	imgStorage          objstore.Storage
	accessTokenDuration time.Duration
}
type USConfig struct {
	userRepo            repository.UserRepository
	patientRepo         repository.PatientRepository
	specialistRepo      repository.SpecialistRepository
	authMaker           auth.Maker
	imgStorage          objstore.Storage
	accessTokenDuration time.Duration
}

func NewUserService(u *USConfig) UserService {
	return &userService{
		userRepo:            u.userRepo,
		patientRepo:         u.patientRepo,
		specialistRepo:      u.specialistRepo,
		authMaker:           u.authMaker,
		imgStorage:          u.imgStorage,
		accessTokenDuration: u.accessTokenDuration,
	}
}

func (s *userService) CreatePatient(ctx context.Context, req model.CreatePatientRequest, userId int64) (database.Patient, error) {
	return s.patientRepo.Create(ctx, repository.CreatePatientParams{
		UserID:      userId,
		Allergies:   req.Allergies,
		DateOfBirth: req.DateOfBirth,
	})
}

func (s *userService) CreateSpecialist(ctx context.Context, req model.CreateSpecialistRequest, userId int64) (database.Specialist, error) {
	return s.specialistRepo.Create(ctx, repository.CreateSpecialistParams{
		Specialization: req.Specialization,
		LicenseNumber:  req.LicenseNumber,
		UserID:         userId,
	})
}

func (s *userService) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.AuthResponse, error) {
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to process password:%v", err)
	}

	user, err := s.userRepo.Create(ctx, repository.CreateUserParams{
		FullName: req.Fullname,
		Email:    req.Email,
		Password: passwordHash,
		UserRole: req.Role,
	})
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("failed to create user:%v", err)
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
	fmt.Println(StringToNullString(imageURL))
	if err != nil {
		return fmt.Errorf("unable to upload the image:%v", err)
	}
	return s.userRepo.UpdateProfilePicture(ctx, StringToNullString(imageURL), userId)
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
	token, err := s.authMaker.Create(user.Email, user.UserID, s.accessTokenDuration)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		AccessToken: token,
		User:        userResponse,
	}, nil
}

func objNameFromURL(imageURL sql.NullString, fileName string) (string, error) {
	// if user doesn't have imageURL - create one
	// otherwise, extract last part of URL to get cloud storage object name
	if !imageURL.Valid {
		objectName := fmt.Sprintf("%s_%d", fileName, time.Now().UnixNano())

		return objectName, nil
	}

	// split off last part of URL, which is the image's storage object ID
	urlPath, err := url.Parse(imageURL.String)
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
