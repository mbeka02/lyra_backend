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
	Create(context.Context, CreatePatientParams) (*database.Patient, error)
	GetPatientIdByUserId(context.Context, int64) (int64, error)
}

type patientRepository struct {
	store *database.Store
}

func NewPatientRepository(store *database.Store) PatientRepository {
	return &patientRepository{
		store,
	}
}

func (r *patientRepository) Create(ctx context.Context, params CreatePatientParams) (*database.Patient, error) {
	var patient database.Patient
	err := r.store.ExecTx(ctx, func(q *database.Queries) error {
		var err error
		// create record
		patient, err = q.CreatePatient(ctx, database.CreatePatientParams{
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
		if err != nil {
			return err
		}
		// mark onboarding as completed
		err = q.CompleteOnboarding(ctx, params.UserID)
		return err
	})
	return &patient, err
}

func (r *patientRepository) GetPatientIdByUserId(ctx context.Context, UserID int64) (int64, error) {
	return r.store.GetPatientIdByUserId(ctx, UserID)
}
