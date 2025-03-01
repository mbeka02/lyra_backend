package repository

import (
	"context"
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateAvailabilityParams struct {
	DoctorID  int64
	DayOfWeek int32
	StartTime time.Time
	EndTime   time.Time
	// IsRecurring bool      `json:"is_recurring"`
}

type AvailabilityRepository interface {
	Create(context.Context, CreateAvailabilityParams) (database.Availability, error)
}

type availabilityRepository struct {
	store *database.Store
}

func NewAvailabilityRepository(store *database.Store) AvailabilityRepository {
	return &availabilityRepository{
		store,
	}
}

func (r *availabilityRepository) Create(ctx context.Context, params CreateAvailabilityParams) (database.Availability, error) {
	return r.store.CreateAvailability(ctx, database.CreateAvailabilityParams{
		DoctorID:  params.DoctorID,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
		DayOfWeek: params.DayOfWeek,
	})
}
