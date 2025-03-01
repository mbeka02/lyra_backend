package service

import (
	"context"
	"errors"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type availabilityService struct {
	availabilityRepo repository.AvailabilityRepository
	doctorRepo       repository.DoctorRepository
}

type AvailabilityService interface {
	CreateAvailability(ctx context.Context, req model.CreateAvailabilityRequest, userId int64) (database.Availability, error)
	GetAvailabilityByDoctor(ctx context.Context, userId int64) ([]database.GetAvailabilityByDoctorRow, error)
}

func NewAvailabilityService(availabilityRepo repository.AvailabilityRepository, doctorRepo repository.DoctorRepository) AvailabilityService {
	return &availabilityService{
		availabilityRepo,
		doctorRepo,
	}
}

func (s *availabilityService) CreateAvailability(ctx context.Context, req model.CreateAvailabilityRequest, userId int64) (database.Availability, error) {
	doctorId, err := s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
	if err != nil {
		return database.Availability{}, errors.New("unable to get the user details of this account")
	}
	return s.availabilityRepo.Create(ctx, repository.CreateAvailabilityParams{
		DoctorID:        doctorId,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DayOfWeek:       req.DayOfWeek,
		IntervalMinutes: req.IntervalMinutes,
	})
}

func (s *availabilityService) GetAvailabilityByDoctor(ctx context.Context, userId int64) ([]database.GetAvailabilityByDoctorRow, error) {
	doctorId, err := s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.New("unable to get the details of this account")
	}

	return s.availabilityRepo.GetByDoctor(ctx, doctorId)
}
