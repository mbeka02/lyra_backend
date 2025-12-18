package repository

import (
	"context"

	"github.com/mbeka02/lyra_backend/internal/database"
)

type paymentRepository struct {
	store *database.Store
}
type UpdatePaymentAndAppointmentStatusParams struct {
	PaymentStatus     string
	AppointmentStatus string
	Reference         string
}
type PaymentRepository interface {
	UpdatePaymentAndAppointmentStatus(ctx context.Context, params UpdatePaymentAndAppointmentStatusParams) error
	GetPaymentByReference(ctx context.Context, reference string) (*database.Payment, error)
}

func NewPaymentRepository(store *database.Store) PaymentRepository {
	return &paymentRepository{store}
}

func (r *paymentRepository) GetPaymentByReference(ctx context.Context, reference string) (*database.Payment, error) {
	payment, err := r.store.GetPaymentByReference(ctx, reference)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdatePaymentAndAppointmentStatus(ctx context.Context, params UpdatePaymentAndAppointmentStatusParams) error {
	// UPDATES THE PAYMENT AND APPOINTMENT STATUS AT THE SAME TIME
	return r.store.ExecTx(ctx, func(q *database.Queries) error {
		var (
			err     error
			payment database.Payment
		)
		// get the payment  record
		payment, err = q.GetPaymentByReference(ctx, params.Reference)

		// update the payment status
		err = q.UpdatePaymentStatus(ctx, database.UpdatePaymentStatusParams{
			CurrentStatus: database.PaymentStatus(params.PaymentStatus),
			Reference:     params.Reference,
		})
		if err != nil {
			return err
		}
		// update the appointment status
		err = q.UpdateAppointmentStatus(ctx, database.UpdateAppointmentStatusParams{
			AppointmentID: payment.AppointmentID,
			CurrentStatus: database.AppointmentStatus(params.AppointmentStatus),
		})

		return err
	})
}
