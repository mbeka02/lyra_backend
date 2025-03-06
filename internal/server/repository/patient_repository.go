package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreatePatientParams struct {
	Allergies             string
	UserID                int64
	CurrentMedication     string
	PastMedicalHistory    string
	FamilyMedicalHistory  string
	InsuranceProvider     string
	InsurancePolicyNumber string
	Address               string
	EmergencyContactName  string
	EmergencyContactPhone string
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
		UserID:                params.UserID,
		Allergies:             params.Allergies,
		CurrentMedication:     params.CurrentMedication,
		PastMedicalHistory:    params.PastMedicalHistory,
		FamilyMedicalHistory:  params.FamilyMedicalHistory,
		InsurancePolicyNumber: params.InsurancePolicyNumber,
		InsuranceProvider:     params.InsuranceProvider,
		Address:               params.Address,
		EmergencyContactName:  params.EmergencyContactName,
		EmergencyContactPhone: params.EmergencyContactPhone,
	})
}
