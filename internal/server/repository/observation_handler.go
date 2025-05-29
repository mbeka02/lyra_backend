// internal/server/repository/ehr_observation_repository.go
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mbeka02/lyra_backend/internal/database" // Your sqlc generated package
)

type ObservationRepository interface {
	Create(ctx context.Context, params database.CreateObservationParams) (database.Observation, error)
	GetByID(ctx context.Context, id uuid.UUID, patientID int64) (database.Observation, error)
	ListByPatientID(ctx context.Context, patientID int64) ([]database.Observation, error)
	Update(ctx context.Context, params database.UpdateObservationParams) (database.Observation, error)
	Delete(ctx context.Context, id uuid.UUID, patientID int64) error
}

type sqlObservationRepository struct {
	store *database.Store
}

func NewSQLObservationRepository(store *database.Store) ObservationRepository {
	return &sqlObservationRepository{store: store}
}

func (r *sqlObservationRepository) Create(ctx context.Context, params database.CreateObservationParams) (database.Observation, error) {
	return r.store.CreateObservation(ctx, params)
}

func (r *sqlObservationRepository) GetByID(ctx context.Context, id uuid.UUID, patientID int64) (database.Observation, error) {
	return r.store.GetObservationByID(ctx, database.GetObservationByIDParams{
		ID:        id,
		PatientID: patientID,
	})
}

func (r *sqlObservationRepository) ListByPatientID(ctx context.Context, patientID int64) ([]database.Observation, error) {
	return r.store.ListObservationsByPatient(ctx, patientID)
}

func (r *sqlObservationRepository) Update(ctx context.Context, params database.UpdateObservationParams) (database.Observation, error) {
	return r.store.UpdateObservation(ctx, params)
}

func (r *sqlObservationRepository) Delete(ctx context.Context, id uuid.UUID, patientID int64) error {
	return r.store.DeleteObservation(ctx, database.DeleteObservationParams{
		ID:        id,
		PatientID: patientID,
	})
}
