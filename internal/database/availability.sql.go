// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: availability.sql

package database

import (
	"context"
	"time"
)

const createAvailability = `-- name: CreateAvailability :one
INSERT INTO availability (
  doctor_id, day_of_week, start_time, end_time, is_recurring,interval_minutes
) VALUES (
  $1, $2, $3, $4, $5,$6
) RETURNING availability_id, doctor_id, start_time, end_time, is_recurring, created_at, updated_at, day_of_week, interval_minutes
`

type CreateAvailabilityParams struct {
	DoctorID        int64  `json:"doctor_id"`
	DayOfWeek       int32  `json:"day_of_week"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	IsRecurring     bool   `json:"is_recurring"`
	IntervalMinutes int32  `json:"interval_minutes"`
}

func (q *Queries) CreateAvailability(ctx context.Context, arg CreateAvailabilityParams) (Availability, error) {
	row := q.db.QueryRowContext(ctx, createAvailability,
		arg.DoctorID,
		arg.DayOfWeek,
		arg.StartTime,
		arg.EndTime,
		arg.IsRecurring,
		arg.IntervalMinutes,
	)
	var i Availability
	err := row.Scan(
		&i.AvailabilityID,
		&i.DoctorID,
		&i.StartTime,
		&i.EndTime,
		&i.IsRecurring,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DayOfWeek,
		&i.IntervalMinutes,
	)
	return i, err
}

const deleteAvailabityByDay = `-- name: DeleteAvailabityByDay :exec
DELETE  FROM availability WHERE day_of_week=$1 AND doctor_id=$2
`

type DeleteAvailabityByDayParams struct {
	DayOfWeek int32 `json:"day_of_week"`
	DoctorID  int64 `json:"doctor_id"`
}

func (q *Queries) DeleteAvailabityByDay(ctx context.Context, arg DeleteAvailabityByDayParams) error {
	_, err := q.db.ExecContext(ctx, deleteAvailabityByDay, arg.DayOfWeek, arg.DoctorID)
	return err
}

const deleteAvailabityById = `-- name: DeleteAvailabityById :exec
DELETE FROM availability WHERE availability_id=$1 AND doctor_id=$2
`

type DeleteAvailabityByIdParams struct {
	AvailabilityID int64 `json:"availability_id"`
	DoctorID       int64 `json:"doctor_id"`
}

func (q *Queries) DeleteAvailabityById(ctx context.Context, arg DeleteAvailabityByIdParams) error {
	_, err := q.db.ExecContext(ctx, deleteAvailabityById, arg.AvailabilityID, arg.DoctorID)
	return err
}

const getAppointmentSlots = `-- name: GetAppointmentSlots :many
WITH time_slots AS (
  SELECT 
    a.doctor_id,
    slot_time::time AS slot_start_time,
    (slot_time + (a.interval_minutes * interval '1 minute'))::time AS slot_end_time
  FROM availability a,
  LATERAL generate_series(
    ($3::date + a.start_time)::timestamp,
    ($3::date + a.end_time - (a.interval_minutes * interval '1 minute'))::timestamp,
    (a.interval_minutes * interval '1 minute')
  ) AS slot_time
  WHERE a.doctor_id = $1
  AND a.day_of_week = $2
)
SELECT
  ts.slot_start_time,
  ts.slot_end_time,
  CASE 
    WHEN appt.appointment_id IS NOT NULL THEN 'booked'
    ELSE 'available'
  END AS slot_status
FROM time_slots ts
LEFT JOIN appointments appt
  ON appt.doctor_id = ts.doctor_id
  AND appt.start_time::time = ts.slot_start_time
  AND appt.start_time::date = $3::date
`

type GetAppointmentSlotsParams struct {
	DoctorID  int64     `json:"doctor_id"`
	DayOfWeek int32     `json:"day_of_week"`
	Column3   time.Time `json:"column_3"`
}

type GetAppointmentSlotsRow struct {
	SlotStartTime string `json:"slot_start_time"`
	SlotEndTime   string `json:"slot_end_time"`
	SlotStatus    string `json:"slot_status"`
}

func (q *Queries) GetAppointmentSlots(ctx context.Context, arg GetAppointmentSlotsParams) ([]GetAppointmentSlotsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAppointmentSlots, arg.DoctorID, arg.DayOfWeek, arg.Column3)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAppointmentSlotsRow
	for rows.Next() {
		var i GetAppointmentSlotsRow
		if err := rows.Scan(&i.SlotStartTime, &i.SlotEndTime, &i.SlotStatus); err != nil {
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

const getAvailabilityByDoctor = `-- name: GetAvailabilityByDoctor :many
SELECT availability_id, doctor_id, start_time, end_time, is_recurring, created_at, updated_at, day_of_week, interval_minutes FROM availability WHERE doctor_id=$1
`

func (q *Queries) GetAvailabilityByDoctor(ctx context.Context, doctorID int64) ([]Availability, error) {
	rows, err := q.db.QueryContext(ctx, getAvailabilityByDoctor, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Availability
	for rows.Next() {
		var i Availability
		if err := rows.Scan(
			&i.AvailabilityID,
			&i.DoctorID,
			&i.StartTime,
			&i.EndTime,
			&i.IsRecurring,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DayOfWeek,
			&i.IntervalMinutes,
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

const getAvailabilityByDoctorAndDay = `-- name: GetAvailabilityByDoctorAndDay :many
SELECT availability_id, doctor_id, start_time, end_time, is_recurring, created_at, updated_at, day_of_week, interval_minutes FROM availability WHERE doctor_id=$1 AND day_of_week=$2
`

type GetAvailabilityByDoctorAndDayParams struct {
	DoctorID  int64 `json:"doctor_id"`
	DayOfWeek int32 `json:"day_of_week"`
}

func (q *Queries) GetAvailabilityByDoctorAndDay(ctx context.Context, arg GetAvailabilityByDoctorAndDayParams) ([]Availability, error) {
	rows, err := q.db.QueryContext(ctx, getAvailabilityByDoctorAndDay, arg.DoctorID, arg.DayOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Availability
	for rows.Next() {
		var i Availability
		if err := rows.Scan(
			&i.AvailabilityID,
			&i.DoctorID,
			&i.StartTime,
			&i.EndTime,
			&i.IsRecurring,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DayOfWeek,
			&i.IntervalMinutes,
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
