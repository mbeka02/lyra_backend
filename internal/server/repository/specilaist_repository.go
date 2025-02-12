package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateSpecialistParams struct {
	Specialization string
	LicenseNumber  string
	UserID         int64
}
type SpecialistRepository interface {
	Create(context.Context, CreateSpecialistParams) (database.Specialist, error)
}

type specialistRepository struct {
	store *database.Store
}

func NewSpecialistRepository(store *database.Store) SpecialistRepository {
	return &specialistRepository{
		store,
	}
}

func (s *specialistRepository) Create(ctx context.Context, params CreateSpecialistParams) (database.Specialist, error) {
	return s.store.CreateSpecialist(ctx, database.CreateSpecialistParams{
		UserID:         params.UserID,
		LicenseNumber:  params.LicenseNumber,
		Specialization: params.Specialization,
	})
}
