package repository

import (
	"context"
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreatePatientParams struct {
	Allergies   string
	DateOfBirth time.Time
	UserID      int64
}
type PatientRepository interface {
	Create(context.Context, CreatePatientParams) (database.Patient, error)
}

type patientRepository struct {
	store *database.Store
}

func NewPatientRepository(store *database.Store) PatientRepository {
	return &patientRepository{
		store,
	}
}

func (p *patientRepository) Create(ctx context.Context, params CreatePatientParams) (database.Patient, error) {
	return p.store.CreatePatient(ctx, database.CreatePatientParams{
		UserID:      params.UserID,
		Allergies:   params.Allergies,
		DateOfBirth: params.DateOfBirth,
	})
}
