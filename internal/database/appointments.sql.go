// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: appointments.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createAppointment = `-- name: CreateAppointment :one
INSERT INTO appointments(patient_id,doctor_id,start_time,end_time, reason) VALUES ($1,$2,$3,$4,$5) RETURNING appointment_id, patient_id, doctor_id, current_status, reason, notes, start_time, end_time, created_at, updated_at
`

type CreateAppointmentParams struct {
	PatientID int64     `json:"patient_id"`
	DoctorID  int64     `json:"doctor_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Reason    string    `json:"reason"`
}

func (q *Queries) CreateAppointment(ctx context.Context, arg CreateAppointmentParams) (Appointment, error) {
	row := q.db.QueryRowContext(ctx, createAppointment,
		arg.PatientID,
		arg.DoctorID,
		arg.StartTime,
		arg.EndTime,
		arg.Reason,
	)
	var i Appointment
	err := row.Scan(
		&i.AppointmentID,
		&i.PatientID,
		&i.DoctorID,
		&i.CurrentStatus,
		&i.Reason,
		&i.Notes,
		&i.StartTime,
		&i.EndTime,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAppointment = `-- name: DeleteAppointment :exec
DELETE FROM appointments WHERE appointment_id=$1
`

func (q *Queries) DeleteAppointment(ctx context.Context, appointmentID int64) error {
	_, err := q.db.ExecContext(ctx, deleteAppointment, appointmentID)
	return err
}

const getPatientAppointments = `-- name: GetPatientAppointments :many
SELECT
a.appointment_id, a.patient_id, a.doctor_id, a.current_status, a.reason, a.notes, a.start_time, a.end_time, a.created_at, a.updated_at,
d.specialization,
u.profile_image_url AS doctor_profile_image_url
FROM 
appointments a
JOIN 
doctors d ON a.doctor_id = d.doctor_id
JOIN 
users u ON d.user_id = u.user_id
WHERE a.patient_id=$1
AND (a.current_status = $2::appointment_status OR TRIM($2::text)='')
AND DATE(a.start_time) BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '1 day'* $3::integer
`

type GetPatientAppointmentsParams struct {
	PatientID   int64             `json:"patient_id"`
	Status      AppointmentStatus `json:"status"`
	SetInterval int32             `json:"set_interval"`
}

type GetPatientAppointmentsRow struct {
	AppointmentID         int64             `json:"appointment_id"`
	PatientID             int64             `json:"patient_id"`
	DoctorID              int64             `json:"doctor_id"`
	CurrentStatus         AppointmentStatus `json:"current_status"`
	Reason                string            `json:"reason"`
	Notes                 sql.NullString    `json:"notes"`
	StartTime             time.Time         `json:"start_time"`
	EndTime               time.Time         `json:"end_time"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             sql.NullTime      `json:"updated_at"`
	Specialization        string            `json:"specialization"`
	DoctorProfileImageUrl string            `json:"doctor_profile_image_url"`
}

func (q *Queries) GetPatientAppointments(ctx context.Context, arg GetPatientAppointmentsParams) ([]GetPatientAppointmentsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPatientAppointments, arg.PatientID, arg.Status, arg.SetInterval)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPatientAppointmentsRow
	for rows.Next() {
		var i GetPatientAppointmentsRow
		if err := rows.Scan(
			&i.AppointmentID,
			&i.PatientID,
			&i.DoctorID,
			&i.CurrentStatus,
			&i.Reason,
			&i.Notes,
			&i.StartTime,
			&i.EndTime,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Specialization,
			&i.DoctorProfileImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAppointmentStatus = `-- name: UpdateAppointmentStatus :exec
UPDATE appointments SET current_status=$1 WHERE appointment_id=$2
`

type UpdateAppointmentStatusParams struct {
	CurrentStatus AppointmentStatus `json:"current_status"`
	AppointmentID int64             `json:"appointment_id"`
}

func (q *Queries) UpdateAppointmentStatus(ctx context.Context, arg UpdateAppointmentStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateAppointmentStatus, arg.CurrentStatus, arg.AppointmentID)
	return err
}
