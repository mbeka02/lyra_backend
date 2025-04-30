package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type CreateAppointmentParams struct {
	DoctorID  int64
	StartTime time.Time
	EndTime   time.Time
	Reason    string
}
type CreateAppointmentWithPaymentParams struct {
	DoctorID  int64
	PatientID int64
	StartTime time.Time
	EndTime   time.Time
	Reason    string
	Reference string // payment reference
	Amount    string // this will be cast to a postgres numeric
}
type GetPatientAppointmentsParams struct {
	PatientID int64
	Interval  int32
	Status    string
}
type GetDoctorAppointmentsParams struct {
	DoctorID int64
	Interval int32
	Status   string
}
type CreateAppointmentWithPaymentTxResults struct {
	Appointment database.Appointment `json:"appointment"`
	Payment     database.Payment     `json:"payment"`
}
type GetAppointmentIDsParams struct {
	ID   int64
	Role string
}
type UpdateAppointmentStatusParams struct {
	AppointmentID int64
	Status        string
}
type CheckAppointmentExistsParams struct {
	PatientID int64
	DoctorID  int64
}
type AppointmentRepository interface {
	CreateAppointmentWithPayment(ctx context.Context, params CreateAppointmentWithPaymentParams) (*CreateAppointmentWithPaymentTxResults, error)
	GetPatientAppointments(ctx context.Context, params GetPatientAppointmentsParams) ([]database.GetPatientAppointmentsRow, error)
	GetDoctorAppointments(ctx context.Context, params GetDoctorAppointmentsParams) ([]database.GetDoctorAppointmentsRow, error)
	GetAppointmentIDs(ctx context.Context, params GetAppointmentIDsParams) ([]int64, error)
	UpdateAppointmentStatus(ctx context.Context, params UpdateAppointmentStatusParams) error
	CheckAppointmentExists(ctx context.Context, params CheckAppointmentExistsParams) (bool, error)
}

type appointmentRepository struct {
	store *database.Store
}

func NewAppointmentRepository(store *database.Store) AppointmentRepository {
	return &appointmentRepository{
		store,
	}
}

func (r *appointmentRepository) CheckAppointmentExists(ctx context.Context, params CheckAppointmentExistsParams) (bool, error) {
	exists, err := r.store.CheckSpecialistPatientAppointmentExists(ctx,
		database.CheckSpecialistPatientAppointmentExistsParams{
			DoctorID:  params.DoctorID,
			PatientID: params.PatientID,
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to check appointment existence: %w", err)
	}
	return exists, nil
}

func (r *appointmentRepository) UpdateAppointmentStatus(ctx context.Context, params UpdateAppointmentStatusParams) error {
	return r.store.UpdateAppointmentStatus(ctx, database.UpdateAppointmentStatusParams{
		CurrentStatus: database.AppointmentStatus(params.Status),
		AppointmentID: params.AppointmentID,
	})
}

func (r *appointmentRepository) GetAppointmentIDs(ctx context.Context, params GetAppointmentIDsParams) ([]int64, error) {
	return r.store.GetAppointmentIDs(ctx, database.GetAppointmentIDsParams{
		ID:   params.ID,
		Role: params.Role,
	})
}

func (r *appointmentRepository) GetDoctorAppointments(ctx context.Context, params GetDoctorAppointmentsParams) ([]database.GetDoctorAppointmentsRow, error) {
	return r.store.GetDoctorAppointments(ctx, database.GetDoctorAppointmentsParams{
		DoctorID:    params.DoctorID,
		SetInterval: params.Interval,
		Status:      params.Status,
	})
}

func (r *appointmentRepository) GetPatientAppointments(ctx context.Context, params GetPatientAppointmentsParams) ([]database.GetPatientAppointmentsRow, error) {
	return r.store.GetPatientAppointments(ctx, database.GetPatientAppointmentsParams{
		PatientID:   params.PatientID,
		SetInterval: params.Interval,
		Status:      params.Status,
	})
}

func (r *appointmentRepository) CreateAppointmentWithPayment(ctx context.Context, params CreateAppointmentWithPaymentParams) (*CreateAppointmentWithPaymentTxResults, error) {
	var result CreateAppointmentWithPaymentTxResults
	err := r.store.ExecTx(ctx, func(q *database.Queries) error {
		var err error
		// create appointment record
		result.Appointment, err = q.CreateAppointment(ctx, database.CreateAppointmentParams{
			DoctorID:  params.DoctorID,
			PatientID: params.PatientID,
			StartTime: params.StartTime,
			EndTime:   params.EndTime,
			Reason:    params.Reason,
		})
		if err != nil {
			return err
		}
		// create payment record
		result.Payment, err = q.CreatePayment(ctx, database.CreatePaymentParams{
			AppointmentID: result.Appointment.AppointmentID,
			DoctorID:      params.DoctorID,
			PatientID:     params.PatientID,
			Reference:     params.Reference,
			Amount:        params.Amount,
		})

		return err
	})
	return &result, err
}
