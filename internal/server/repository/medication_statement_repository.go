package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/database"
)

type MedicationStatementRepository interface {
	Create(ctx context.Context, params database.CreateMedicationStatementParams) (database.MedicationStatement, error)
	GetByID(ctx context.Context, id uuid.UUID, patientID int64) (database.MedicationStatement, error)
	ListByPatientID(ctx context.Context, patientID int64) ([]database.MedicationStatement, error)
	Update(ctx context.Context, params database.UpdateMedicationStatementParams) (database.MedicationStatement, error)
	Delete(ctx context.Context, id uuid.UUID, patientID int64) error
}

type sqlMedicationStatementRepository struct {
	store *database.Store
}

func NewSQLMedicationStatementRepository(store *database.Store) MedicationStatementRepository {
	return &sqlMedicationStatementRepository{store: store}
}

func (r *sqlMedicationStatementRepository) Create(ctx context.Context, params database.CreateMedicationStatementParams) (database.MedicationStatement, error) {
	return r.store.CreateMedicationStatement(ctx, params)
}

func (r *sqlMedicationStatementRepository) GetByID(ctx context.Context, id uuid.UUID, patientID int64) (database.MedicationStatement, error) {
	return r.store.GetMedicationStatementByID(ctx, database.GetMedicationStatementByIDParams{
		ID:        id,
		PatientID: patientID,
	})
}

func (r *sqlMedicationStatementRepository) ListByPatientID(ctx context.Context, patientID int64) ([]database.MedicationStatement, error) {
	return r.store.ListMedicationStatementsByPatient(ctx, patientID)
}

func (r *sqlMedicationStatementRepository) Update(ctx context.Context, params database.UpdateMedicationStatementParams) (database.MedicationStatement, error) {
	return r.store.UpdateMedicationStatement(ctx, params)
}

func (r *sqlMedicationStatementRepository) Delete(ctx context.Context, id uuid.UUID, patientID int64) error {
	return r.store.DeleteMedicationStatement(ctx, database.DeleteMedicationStatementParams{
		ID:        id,
		PatientID: patientID,
	})
}
