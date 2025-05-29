// internal/server/service/ehr_allergy_service.go
package service

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type AllergyService interface {
	CreateAllergy(ctx context.Context, req model.CreateAllergyIntoleranceRequest, actingUserID int64, forPatientID int64) (database.AllergyIntolerance, error)
	GetAllergy(ctx context.Context, allergyID uuid.UUID, actingUserID int64, forPatientID int64) (database.AllergyIntolerance, error)
	ListAllergiesForPatient(ctx context.Context, actingUserID int64, forPatientID int64) ([]database.AllergyIntolerance, error)
	UpdateAllergy(ctx context.Context, allergyID uuid.UUID, req model.UpdateAllergyIntoleranceRequest, actingUserID int64, forPatientID int64) (database.AllergyIntolerance, error)
	DeleteAllergy(ctx context.Context, allergyID uuid.UUID, actingUserID int64, forPatientID int64) error
}

type allergyService struct {
	allergyRepo repository.AllergyIntoleranceRepository
	// patientRepo    repository.PatientRepository // TODO: validate patient existence
}

func NewAllergyService(allergyRepo repository.AllergyIntoleranceRepository) AllergyService {
	return &allergyService{
		allergyRepo: allergyRepo,
	}
}

func (s *allergyService) CreateAllergy(ctx context.Context, req model.CreateAllergyIntoleranceRequest, actingUserID int64, forPatientID int64) (database.AllergyIntolerance, error) {
	// TODO: Authorization: Can actingUserID create allergy for forPatientID?

	params := database.CreateAllergyIntoleranceParams{
		PatientID:                 forPatientID,
		ClinicalStatusCode:        req.ClinicalStatusCode,
		ClinicalStatusDisplay:     ToNullString(req.ClinicalStatusDisplay),
		CodeSystem:                ToNullString(req.CodeSystem),
		CodeCode:                  req.CodeCode,
		CodeDisplay:               req.CodeDisplay,
		Criticality:               ToNullString(req.Criticality),
		ReactionManifestationText: ToNullString(req.ReactionManifestationText),
	}
	return s.allergyRepo.Create(ctx, params)
}

func (s *allergyService) GetAllergy(ctx context.Context, allergyID uuid.UUID, actingUserID int64, forPatientID int64) (database.AllergyIntolerance, error) {
	// TODO: Authorization: Can actingUserID view allergy for forPatientID?
	return s.allergyRepo.GetByID(ctx, allergyID, forPatientID)
}

func (s *allergyService) ListAllergiesForPatient(ctx context.Context, actingUserID int64, forPatientID int64) ([]database.AllergyIntolerance, error) {
	// TODO: Authorization: Can actingUserID list allergies for forPatientID?
	return s.allergyRepo.ListByPatientID(ctx, forPatientID)
}

func (s *allergyService) UpdateAllergy(ctx context.Context, allergyID uuid.UUID, req model.UpdateAllergyIntoleranceRequest, actingUserID int64, forPatientID int64) (database.AllergyIntolerance, error) {
	// TODO: Authorization: Can actingUserID update allergy for forPatientID?

	params := database.UpdateAllergyIntoleranceParams{
		ID:                        allergyID,
		PatientID:                 forPatientID, // Used in WHERE clause for safety
		ClinicalStatusCode:        req.ClinicalStatusCode,
		ClinicalStatusDisplay:     ToNullString(req.ClinicalStatusDisplay),
		CodeSystem:                ToNullString(req.CodeSystem),
		CodeCode:                  req.CodeCode,
		CodeDisplay:               req.CodeDisplay,
		Criticality:               ToNullString(req.Criticality),
		ReactionManifestationText: ToNullString(req.ReactionManifestationText),
	}
	return s.allergyRepo.Update(ctx, params)
}

func (s *allergyService) DeleteAllergy(ctx context.Context, allergyID uuid.UUID, actingUserID int64, forPatientID int64) error {
	// TODO: Authorization: Can actingUserID delete allergy for forPatientID?

	return s.allergyRepo.Delete(ctx, allergyID, forPatientID)
}

// Helper for nullable strings (put in a common utils package or keep here for now)
func ToNullString(s *string) sql.NullString {
	if s == nil || *s == "" { // Treat empty string from JSON as NULL in DB
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}
