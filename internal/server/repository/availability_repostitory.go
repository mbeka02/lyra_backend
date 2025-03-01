package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateAvailabilityParams struct {
	DoctorID        int64
	DayOfWeek       int32
	StartTime       string
	EndTime         string
	IntervalMinutes int32

	// IsRecurring bool      `json:"is_recurring"`
}

type AvailabilityRepository interface {
	Create(context.Context, CreateAvailabilityParams) (database.Availability, error)
	GetByDoctor(ctx context.Context, doctorId int64) ([]database.Availability, error)
	DeleteById(ctx context.Context, availabilityId int64, doctorId int64) error
	DeleteByDay(ctx context.Context, dayOfWeek int32, doctorId int64) error
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
		DoctorID:        params.DoctorID,
		StartTime:       params.StartTime,
		EndTime:         params.EndTime,
		DayOfWeek:       params.DayOfWeek,
		IntervalMinutes: params.IntervalMinutes,
	})
}

func (r *availabilityRepository) GetByDoctor(ctx context.Context, DoctorID int64) ([]database.Availability, error) {
	return r.store.GetAvailabilityByDoctor(ctx, DoctorID)
}

func (r *availabilityRepository) DeleteById(ctx context.Context, availabilityId int64, doctorId int64) error {
	return r.store.DeleteAvailabityById(ctx, database.DeleteAvailabityByIdParams{
		AvailabilityID: availabilityId,
		DoctorID:       doctorId,
	})
}

func (r *availabilityRepository) DeleteByDay(ctx context.Context, dayOfWeek int32, doctorId int64) error {
	return r.store.DeleteAvailabityByDay(ctx, database.DeleteAvailabityByDayParams{
		DayOfWeek: dayOfWeek,
		DoctorID:  doctorId,
	})
}
