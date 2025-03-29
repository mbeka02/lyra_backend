// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: payments.sql

package database

import (
	"context"
)

const createPayment = `-- name: CreatePayment :one
INSERT INTO payments (
  reference,
  amount,
  patient_id,
  doctor_id,
  appointment_id
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING payment_id, reference, current_status, amount, metadata, payment_method, currency, appointment_id, patient_id, doctor_id, created_at, updated_at, completed_at
`

type CreatePaymentParams struct {
	Reference     string `json:"reference"`
	Amount        string `json:"amount"`
	PatientID     int64  `json:"patient_id"`
	DoctorID      int64  `json:"doctor_id"`
	AppointmentID int64  `json:"appointment_id"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRowContext(ctx, createPayment,
		arg.Reference,
		arg.Amount,
		arg.PatientID,
		arg.DoctorID,
		arg.AppointmentID,
	)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.Reference,
		&i.CurrentStatus,
		&i.Amount,
		&i.Metadata,
		&i.PaymentMethod,
		&i.Currency,
		&i.AppointmentID,
		&i.PatientID,
		&i.DoctorID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CompletedAt,
	)
	return i, err
}

const getPaymentByReference = `-- name: GetPaymentByReference :one
SELECT payment_id, reference, current_status, amount, metadata, payment_method, currency, appointment_id, patient_id, doctor_id, created_at, updated_at, completed_at FROM payments WHERE reference = $1 LIMIT 1
`

func (q *Queries) GetPaymentByReference(ctx context.Context, reference string) (Payment, error) {
	row := q.db.QueryRowContext(ctx, getPaymentByReference, reference)
	var i Payment
	err := row.Scan(
		&i.PaymentID,
		&i.Reference,
		&i.CurrentStatus,
		&i.Amount,
		&i.Metadata,
		&i.PaymentMethod,
		&i.Currency,
		&i.AppointmentID,
		&i.PatientID,
		&i.DoctorID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CompletedAt,
	)
	return i, err
}

const updatePaymentStatus = `-- name: UpdatePaymentStatus :exec
UPDATE payments
SET 
  current_status = $1,
  completed_at = CASE WHEN $1 = 'completed'::payment_status THEN NOW() ELSE completed_at END,
  updated_at = NOW()
WHERE reference = $2
`

type UpdatePaymentStatusParams struct {
	CurrentStatus PaymentStatus `json:"current_status"`
	Reference     string        `json:"reference"`
}

func (q *Queries) UpdatePaymentStatus(ctx context.Context, arg UpdatePaymentStatusParams) error {
	_, err := q.db.ExecContext(ctx, updatePaymentStatus, arg.CurrentStatus, arg.Reference)
	return err
}
