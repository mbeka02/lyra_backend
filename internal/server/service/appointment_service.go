package service

import (
	"context"
	"errors"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	patientRepo     repository.PatientRepository
}

type AppointmentService interface {
	CreateAppointment(ctx context.Context, req model.CreateAppointmentRequest, userId int64) (database.Appointment, error)
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, patientRepo repository.PatientRepository) AppointmentService {
	return &appointmentService{
		appointmentRepo,
		patientRepo,
	}
}

func (s *appointmentService) CreateAppointment(ctx context.Context, req model.CreateAppointmentRequest, userId int64) (database.Appointment, error) {
	patientId, err := s.patientRepo.GetPatientIdByUserId(ctx, userId)
	if err != nil {
		return database.Appointment{}, errors.New("unable to get the user details of this account")
	}

	return s.appointmentRepo.Create(ctx, repository.CreateAppointmentParams{
		DoctorID:  req.DoctorID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Reason:    req.Reason,
	}, patientId)
}
