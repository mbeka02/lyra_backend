// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: doctors.sql

package database

import (
	"context"
)

const createDoctor = `-- name: CreateDoctor :one
INSERT INTO doctors(user_id,specialization,license_number,description) VALUES ($1,$2,$3,$4)RETURNING doctor_id, user_id, description, specialization, license_number, created_at, updated_at
`

type CreateDoctorParams struct {
	UserID         int64
	Specialization string
	LicenseNumber  string
	Description    string
}

func (q *Queries) CreateDoctor(ctx context.Context, arg CreateDoctorParams) (Doctor, error) {
	row := q.db.QueryRowContext(ctx, createDoctor,
		arg.UserID,
		arg.Specialization,
		arg.LicenseNumber,
		arg.Description,
	)
	var i Doctor
	err := row.Scan(
		&i.DoctorID,
		&i.UserID,
		&i.Description,
		&i.Specialization,
		&i.LicenseNumber,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getDoctors = `-- name: GetDoctors :many
SELECT full_name,specialization,doctor_id,profile_image_url,description FROM doctors INNER JOIN users ON doctors.user_id=users.user_id LIMIT $1 OFFSET $2
`

type GetDoctorsParams struct {
	Limit  int32
	Offset int32
}

type GetDoctorsRow struct {
	FullName        string
	Specialization  string
	DoctorID        int64
	ProfileImageUrl string
	Description     string
}

func (q *Queries) GetDoctors(ctx context.Context, arg GetDoctorsParams) ([]GetDoctorsRow, error) {
	rows, err := q.db.QueryContext(ctx, getDoctors, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDoctorsRow
	for rows.Next() {
		var i GetDoctorsRow
		if err := rows.Scan(
			&i.FullName,
			&i.Specialization,
			&i.DoctorID,
			&i.ProfileImageUrl,
			&i.Description,
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
