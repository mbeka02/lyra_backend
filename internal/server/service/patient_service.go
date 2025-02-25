package service

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type patientService struct {
	patientRepo repository.PatientRepository
}

type PatientService interface {
	CreatePatient(ctx context.Context, req model.CreatePatientRequest, userId int64) (database.Patient, error)
}

func NewPatientService(patientRepo repository.PatientRepository) PatientService {
	return &patientService{
		patientRepo,
	}
}

func (s *patientService) CreatePatient(ctx context.Context, req model.CreatePatientRequest, userId int64) (database.Patient, error) {
	return s.patientRepo.Create(ctx, repository.CreatePatientParams{
		UserID:    userId,
		Allergies: req.Allergies,
	})
}
