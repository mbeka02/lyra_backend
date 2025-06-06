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
	CreateAvailability(ctx context.Context, req model.CreateAvailabilityRequest, userId int64) (*database.Availability, error)
	GetAvailabilityByDoctor(ctx context.Context, userId int64) ([]database.Availability, error)
	GetSlots(ctx context.Context, req model.GetSlotsRequest) ([]database.GetAppointmentSlotsRow, error)
	DeleteById(ctx context.Context, avavailabilityId int64, userId int64) error
	DeleteByDay(ctx context.Context, dayOfWeek int32, userId int64) error
}

func (s *availabilityService) GetSlots(ctx context.Context, req model.GetSlotsRequest) ([]database.GetAppointmentSlotsRow, error) {
	return s.availabilityRepo.GetSlots(ctx, repository.GetSlotsParams{
		DoctorID:  req.DoctorID,
		DayOfWeek: req.DayOfWeek,
		SlotDate:  req.SlotDate,
	})
}

func (s *availabilityService) DeleteById(ctx context.Context, avavailabilityId int64, userId int64) error {
	doctorId, err := s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
	if err != nil {
		return errors.New("unable to get the user details of this account")
	}
	return s.availabilityRepo.DeleteById(ctx, avavailabilityId, doctorId)
}

func (s *availabilityService) DeleteByDay(ctx context.Context, dayOfWeek int32, userId int64) error {
	doctorId, err := s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
	if err != nil {
		return errors.New("unable to get the user details of this account")
	}
	return s.availabilityRepo.DeleteByDay(ctx, dayOfWeek, doctorId)
}

func NewAvailabilityService(availabilityRepo repository.AvailabilityRepository, doctorRepo repository.DoctorRepository) AvailabilityService {
	return &availabilityService{
		availabilityRepo,
		doctorRepo,
	}
}

func (s *availabilityService) CreateAvailability(ctx context.Context, req model.CreateAvailabilityRequest, userId int64) (*database.Availability, error) {
	doctorId, err := s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.New("unable to get the user details of this account")
	}

	return s.availabilityRepo.Create(ctx, repository.CreateAvailabilityParams{
		DoctorID:        doctorId,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DayOfWeek:       req.DayOfWeek,
		IntervalMinutes: req.IntervalMinutes,
	})
}

func (s *availabilityService) GetAvailabilityByDoctor(ctx context.Context, userId int64) ([]database.Availability, error) {
	doctorId, err := s.doctorRepo.GetDoctorIdByUserId(ctx, userId)
	if err != nil {
		return nil, errors.New("unable to get the details of this account")
	}

	return s.availabilityRepo.GetByDoctor(ctx, doctorId)
}
