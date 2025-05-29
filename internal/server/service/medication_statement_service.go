package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/model"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
)

type MedicationService interface {
	CreateMedication(ctx context.Context, req model.CreateMedicationStatementRequest, actingUserID int64, forPatientID int64) (database.MedicationStatement, error)
	GetMedication(ctx context.Context, medicationID uuid.UUID, actingUserID int64, forPatientID int64) (database.MedicationStatement, error)
	ListMedicationsForPatient(ctx context.Context, actingUserID int64, forPatientID int64) ([]database.MedicationStatement, error)
	UpdateMedication(ctx context.Context, medicationID uuid.UUID, req model.UpdateMedicationStatementRequest, actingUserID int64, forPatientID int64) (database.MedicationStatement, error)
	DeleteMedication(ctx context.Context, medicationID uuid.UUID, actingUserID int64, forPatientID int64) error
}

type medicationService struct {
	medicationRepo repository.MedicationStatementRepository
}

func NewMedicationService(medRepo repository.MedicationStatementRepository) MedicationService {
	return &medicationService{medicationRepo: medRepo}
}

func (s *medicationService) CreateMedication(ctx context.Context, req model.CreateMedicationStatementRequest, actingUserID int64, forPatientID int64) (database.MedicationStatement, error) {
	// TODO: Authorization
	params := database.CreateMedicationStatementParams{
		PatientID:             forPatientID,
		Status:                req.Status,
		MedicationCodeSystem:  ToNullString(req.MedicationCodeSystem),
		MedicationCodeCode:    req.MedicationCodeCode,
		MedicationCodeDisplay: req.MedicationCodeDisplay,
		DosageText:            ToNullString(req.DosageText),
		EffectiveDateTime:     ToNullTime(req.EffectiveDateTime),
	}
	return s.medicationRepo.Create(ctx, params)
}

func (s *medicationService) GetMedication(ctx context.Context, medicationID uuid.UUID, actingUserID int64, forPatientID int64) (database.MedicationStatement, error) {
	// TODO: Authorization
	return s.medicationRepo.GetByID(ctx, medicationID, forPatientID)
}

func (s *medicationService) ListMedicationsForPatient(ctx context.Context, actingUserID int64, forPatientID int64) ([]database.MedicationStatement, error) {
	// TODO: Authorization
	return s.medicationRepo.ListByPatientID(ctx, forPatientID)
}

func (s *medicationService) UpdateMedication(ctx context.Context, medicationID uuid.UUID, req model.UpdateMedicationStatementRequest, actingUserID int64, forPatientID int64) (database.MedicationStatement, error) {
	// TODO: Authorization
	params := database.UpdateMedicationStatementParams{
		ID:                    medicationID,
		PatientID:             forPatientID, // For WHERE clause
		Status:                req.Status,
		MedicationCodeSystem:  ToNullString(req.MedicationCodeSystem),
		MedicationCodeCode:    req.MedicationCodeCode,
		MedicationCodeDisplay: req.MedicationCodeDisplay,
		DosageText:            ToNullString(req.DosageText),
		EffectiveDateTime:     ToNullTime(req.EffectiveDateTime),
	}
	return s.medicationRepo.Update(ctx, params)
}

func (s *medicationService) DeleteMedication(ctx context.Context, medicationID uuid.UUID, actingUserID int64, forPatientID int64) error {
	// TODO: Authorization
	return s.medicationRepo.Delete(ctx, medicationID, forPatientID)
}

func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}
