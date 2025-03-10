package repository

import (
	"context"
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateAppointmentParams struct {
	DoctorID  int64
	StartTime time.Time
	EndTime   time.Time
	Reason    string
}

type AppointmentRepository interface {
	Create(ctx context.Context, params CreateAppointmentParams, PatientID int64) (database.Appointment, error)
}

type appointmentRepository struct {
	store *database.Store
}

func NewAppointmentRepository(store *database.Store) AppointmentRepository {
	return &appointmentRepository{
		store,
	}
}

func (r *appointmentRepository) Create(ctx context.Context, params CreateAppointmentParams, PatientID int64) (database.Appointment, error) {
	return r.store.CreateAppointment(ctx, database.CreateAppointmentParams{
		DoctorID:  params.DoctorID,
		PatientID: PatientID,
		StartTime: params.StartTime,
		EndTime:   params.EndTime,
		Reason:    params.Reason,
	})
}
