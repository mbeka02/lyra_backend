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
type CreatePatientTxResult struct {
	User    database.User
	Patient database.Patient
}
type PatientRepository interface {
	Create(context.Context, CreatePatientParams) (*CreatePatientTxResult, error)
	GetPatientIdByUserId(context.Context, int64) (int64, error)
	UpdateFHIRVersion(ctx context.Context, patientID int64, version string) error
}

type patientRepository struct {
	store *database.Store
}

func NewPatientRepository(store *database.Store) PatientRepository {
	return &patientRepository{
		store,
	}
}

func (r *patientRepository) Create(ctx context.Context, params CreatePatientParams) (*CreatePatientTxResult, error) {
	var result CreatePatientTxResult
	err := r.store.ExecTx(ctx, func(q *database.Queries) error {
		var err error
		// create record
		result.Patient, err = q.CreatePatient(ctx, database.CreatePatientParams{
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
		// get the user record
		result.User, err = q.GetUserById(ctx, params.UserID)
		if err != nil {
			return err
		}
		// mark onboarding as completed
		err = q.CompleteOnboarding(ctx, params.UserID)
		return err
	})
	return &result, err
}

func (r *patientRepository) GetPatientIdByUserId(ctx context.Context, UserID int64) (int64, error) {
	return r.store.GetPatientIdByUserId(ctx, UserID)
}

func (r *patientRepository) UpdateFHIRVersion(ctx context.Context, patientID int64, version string) error {
	return r.store.UpdateFhirVersionId(ctx, database.UpdateFhirVersionIdParams{
		PatientID:   patientID,
		FhirVersion: version,
	})
}
