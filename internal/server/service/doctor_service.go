package service

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type DoctorService interface {
	CreateDoctor(ctx context.Context, req model.CreateDoctorRequest, userId int64) (database.Doctor, error)
	GetDoctors(ctx context.Context, limit int32, offset int32) ([]database.GetDoctorsRow, error)
}
type doctorService struct {
	doctorRepo repository.DoctorRepository
}

func NewDoctorService(doctorRepo repository.DoctorRepository) DoctorService {
	return &doctorService{
		doctorRepo,
	}
}

func (s *doctorService) CreateDoctor(ctx context.Context, req model.CreateDoctorRequest, userId int64) (database.Doctor, error) {
	return s.doctorRepo.Create(ctx, repository.CreateDoctorParams{
		Specialization: req.Specialization,
		LicenseNumber:  req.LicenseNumber,
		Description:    req.Description,
		UserID:         userId,
	})
}

func (s *doctorService) GetDoctors(ctx context.Context, limit int32, offset int32) ([]database.GetDoctorsRow, error) {
	return s.doctorRepo.GetAll(ctx, repository.GetDoctorsParams{
		Limit:  limit,
		Offset: offset,
	})
}
